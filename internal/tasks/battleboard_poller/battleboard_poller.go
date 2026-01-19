package battleboard_poller

import (
	"albionstats/internal/postgres"
	"albionstats/internal/tasks"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"time"
)

type Config struct {
	APIClient      *tasks.Client
	Postgres       *postgres.Postgres
	Logger         *slog.Logger
	Region         string
	PageSize       int
	MaxPages       int
	EventsInterval time.Duration
}

type BattleboardPoller struct {
	apiClient      *tasks.Client
	postgres       *postgres.Postgres
	log            *slog.Logger
	eventsInterval time.Duration
	pageSize       int
	maxPages       int
	region         string
}

func NewBattleboardPoller(cfg Config) (*BattleboardPoller) {
	return &BattleboardPoller{
		apiClient:      cfg.APIClient,
		postgres:       cfg.Postgres,
		log:            cfg.Logger.With("component", "battleboard_poller", "region", cfg.Region),
		eventsInterval: cfg.EventsInterval,
		pageSize:       cfg.PageSize,
		maxPages:       cfg.MaxPages,
		region:         cfg.Region,
	}
}

func (p *BattleboardPoller) Run() {
	p.log.Info("battleboard polling started", "interval", p.eventsInterval, "page_size", p.pageSize, "max_pages", p.maxPages)

	ticker := time.NewTicker(p.eventsInterval)
	defer ticker.Stop()

	p.runBatch() // Run once immediately

	for range ticker.C {
		p.runBatch()
	}
}

func (p *BattleboardPoller) runBatch() {
	var allBattles []tasks.Battle

	// Iterate over max pages to collect all battles
	for page := 0; page < p.maxPages; page++ {
		offset := page * p.pageSize
		battles, err := p.fetchBattlesWithRetry(p.region, offset, p.pageSize)
		if err != nil {
			p.log.Error("failed to fetch battles after retries", "error", err, "page", page, "offset", offset)
			return
		}

		// If no battles returned, we've reached the end
		if len(battles) == 0 {
			break
		}

		// Append battles to allBattles slice
		allBattles = append(allBattles, battles...)
	}

	summaries := p.collectBattleSummaries(allBattles)
	allianceStats := p.collectBattleAllianceStats(allBattles)
	guildStats := p.collectBattleGuildStats(allBattles)
	playerStats := p.collectBattlePlayerStats(allBattles)
	queues := p.collectBattleQueues(allBattles)

	if err := p.postgres.InsertBattleSummaries(summaries); err != nil {
		p.log.Error("failed to insert battle summaries", "error", err)
		return
	}

	if err := p.postgres.InsertBattleAllianceStats(allianceStats); err != nil {
		p.log.Error("failed to insert battle alliance stats", "error", err)
		return
	}

	if err := p.postgres.InsertBattleGuildStats(guildStats); err != nil {
		p.log.Error("failed to insert battle guild stats", "error", err)
		return
	}

	if err := p.postgres.InsertBattlePlayerStats(playerStats); err != nil {
		p.log.Error("failed to insert battle player stats", "error", err)
		return
	}

	if err := p.postgres.InsertBattleQueues(queues); err != nil {
		p.log.Error("failed to insert battle queues", "error", err)
		return
	}

	p.log.Info("battleboard polling completed", "battles", len(allBattles), "summaries", len(summaries),
	 "alliance_stats", len(allianceStats), "guild_stats", len(guildStats), "player_stats", len(playerStats), "queues", len(queues))
}

func (p *BattleboardPoller) fetchBattlesWithRetry(region string, offset, limit int) ([]tasks.Battle, error) {
	const maxRetries = 3
	const baseDelay = time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		battles, err := p.apiClient.FetchBattles(region, offset, limit)
		if err == nil {
			return battles, nil
		}

		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
			p.log.Warn("fetch battles failed, retrying", "error", err, "attempt", attempt+1, "delay", delay)
			time.Sleep(delay)
		} else {
			p.log.Error("fetch battles failed after all retries", "error", err, "attempts", maxRetries)
			return nil, err
		}
	}

	return nil, fmt.Errorf("unexpected error in retry logic")
}

func (p *BattleboardPoller) collectBattleSummaries(battles []tasks.Battle) ([]postgres.BattleSummary){
	summary := make([]postgres.BattleSummary, 0, len(battles))

	for _, battle := range battles {

		guildNumParticipants := make(map[string]int32)
		allianceNumParticipants := make(map[string]int32)
		for _, player := range battle.Players {
			if player.GuildName != nil && *player.GuildName != "" {
				guildNumParticipants[*player.GuildName]++
			}
			if player.AllianceName != nil && *player.AllianceName != "" {
				allianceNumParticipants[*player.AllianceName]++
			}
		}

		allianceSlice := make([]tasks.BattleAlliance, 0, len(battle.Alliances))
		for _, alliance := range battle.Alliances {
			allianceSlice = append(allianceSlice, alliance)
		}
		sort.Slice(allianceSlice, func(i, j int) bool {
			return allianceSlice[i].Kills > allianceSlice[j].Kills
		})
		allianceNames := make([]string, 0, len(battle.Alliances))
		for _, alliance := range allianceSlice {
			name := alliance.Name + " (" + strconv.Itoa(int(allianceNumParticipants[alliance.Name])) + ")"
			allianceNames = append(allianceNames, name)
		}
	
		guildSlice := make([]tasks.BattleGuild, 0, len(battle.Guilds))
		for _, guild := range battle.Guilds {
			guildSlice = append(guildSlice, guild)
		}
		sort.Slice(guildSlice, func(i, j int) bool {
			return guildSlice[i].Kills > guildSlice[j].Kills
		})
		guildNames := make([]string, 0, len(battle.Guilds))
		for _, guild := range guildSlice {
			name := guild.Name + " (" + strconv.Itoa(int(guildNumParticipants[guild.Name])) + ")"
			guildNames = append(guildNames, name)
		}
	
		playerSlice := make([]tasks.BattlePlayer, 0, len(battle.Players))
		for _, player := range battle.Players {
			playerSlice = append(playerSlice, player)
		}
		sort.Slice(playerSlice, func(i, j int) bool {
			return playerSlice[i].Kills > playerSlice[j].Kills
		})
		playerNames := make([]string, 0, len(battle.Players))
		for _, player := range playerSlice {
			playerNames = append(playerNames, player.Name)
		}

		summary = append(summary, postgres.BattleSummary{
			Region:       postgres.Region(p.region),
			BattleID:     battle.ID,
			StartTime:    battle.StartTime,
			EndTime:      battle.EndTime,
			TotalPlayers: int32(len(battle.Players)),
			TotalKills:   int32(battle.TotalKills),
			TotalFame:    battle.TotalFame,
			AllianceNames: allianceNames,
			GuildNames: guildNames,
			PlayerNames: playerNames,
		})
	}

	return summary
}

func (p *BattleboardPoller) collectBattleAllianceStats(battles []tasks.Battle) ([]postgres.BattleAllianceStats){
	allianceStats := make([]postgres.BattleAllianceStats, 0)
	for _, battle := range battles {
		for _, alliance := range battle.Alliances {
			playerCount := 0
			for _, player := range battle.Players {
				if player.AllianceName != nil && *player.AllianceName == alliance.Name {
					playerCount++
				}
			}
	
			allianceStats = append(allianceStats, postgres.BattleAllianceStats{
				Region:       postgres.Region(p.region),
				BattleID:     battle.ID,
				AllianceName: alliance.Name,
				PlayerCount:  int32(playerCount),
				Kills:        int32(alliance.Kills),
				Deaths:       int32(alliance.Deaths),
				KillFame:     int64(alliance.KillFame),
			})
		}
	}
	return allianceStats
}

func (p *BattleboardPoller) collectBattleGuildStats(battles []tasks.Battle) ([]postgres.BattleGuildStats){
	guildStats := make([]postgres.BattleGuildStats, 0)
	for _, battle := range battles {
		for _, guild := range battle.Guilds {
			playerCount := 0
			for _, player := range battle.Players {
				if player.GuildName != nil && *player.GuildName == guild.Name {
					playerCount++
				}
			}
	
			guildStats = append(guildStats, postgres.BattleGuildStats{
				Region:       postgres.Region(p.region),
				BattleID:     battle.ID,
				GuildName:    guild.Name,
				AllianceName: guild.Alliance,
				PlayerCount:  int32(playerCount),
				Kills:        int32(guild.Kills),
				Deaths:       int32(guild.Deaths),
				KillFame:     int64(guild.KillFame),
			})
		}
	}
	return guildStats
}

func (p *BattleboardPoller) collectBattlePlayerStats(battles []tasks.Battle) ([]postgres.BattlePlayerStats){
	playerStats := make([]postgres.BattlePlayerStats, 0)
	for _, battle := range battles {
		for _, player := range battle.Players {
			playerStats = append(playerStats, postgres.BattlePlayerStats{
				Region:       postgres.Region(p.region),
				BattleID:     battle.ID,
				PlayerName:   player.Name,
				GuildName:    player.GuildName,
				AllianceName: player.AllianceName,
				Kills:        int32(player.Kills),
				Deaths:       int32(player.Deaths),
				KillFame:     int64(player.KillFame),
			})
		}
	}
	return playerStats
}

func (p *BattleboardPoller) collectBattleQueues(battles []tasks.Battle) ([]postgres.BattleQueue){
	queues := make([]postgres.BattleQueue, 0, len(battles))

	for _, battle := range battles {
		queues = append(queues, postgres.BattleQueue{
			Region:       postgres.Region(p.region),
			BattleID:     battle.ID,
			TS:           battle.StartTime,
			ErrorCount:   0,
		})
	}
	return queues
}