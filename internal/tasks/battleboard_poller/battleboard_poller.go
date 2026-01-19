package battleboard_poller

import (
	"albionstats/internal/postgres"
	"albionstats/internal/tasks"
	"log/slog"
	"sort"
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
	region         string
}

func NewBattleboardPoller(cfg Config) (*BattleboardPoller) {
	return &BattleboardPoller{
		apiClient:      cfg.APIClient,
		postgres:       cfg.Postgres,
		log:            cfg.Logger.With("component", "battleboard_poller", "region", cfg.Region),
		eventsInterval: cfg.EventsInterval,
		pageSize:       cfg.PageSize,
		region:         cfg.Region,
	}
}

func (p *BattleboardPoller) Run() {
	p.log.Info("battleboard polling started", "interval", p.eventsInterval, "page_size", p.pageSize)

	ticker := time.NewTicker(p.eventsInterval)
	defer ticker.Stop()

	p.runBatch() // Run once immediately

	for range ticker.C {
		p.runBatch()
	}
}

func (p *BattleboardPoller) runBatch() {
	battles, err := p.apiClient.FetchBattles(p.region, 0, p.pageSize)
	if err != nil {
		p.log.Error("failed to fetch battles", "error", err)
		return
	}

	summaries := p.collectBattleSummaries(battles)
	allianceStats := p.collectBattleAllianceStats(battles)
	guildStats := p.collectBattleGuildStats(battles)
	playerStats := p.collectBattlePlayerStats(battles)
	queues := p.collectBattleQueues(battles)

	err = p.postgres.InsertBattleSummaries(summaries)
	if err != nil {
		p.log.Error("failed to insert battle summaries", "error", err)
		return
	}

	err = p.postgres.InsertBattleAllianceStats(allianceStats)
	if err != nil {
		p.log.Error("failed to insert battle alliance stats", "error", err)
		return
	}

	err = p.postgres.InsertBattleGuildStats(guildStats)
	if err != nil {
		p.log.Error("failed to insert battle guild stats", "error", err)
		return
	}

	err = p.postgres.InsertBattlePlayerStats(playerStats)
	if err != nil {
		p.log.Error("failed to insert battle player stats", "error", err)
		return
	}

	err = p.postgres.InsertBattleQueues(queues)
	if err != nil {
		p.log.Error("failed to insert battle queues", "error", err)
		return
	}

	p.log.Info("battleboard polling completed", "battles", len(battles), "summaries", len(summaries),
	 "alliance_stats", len(allianceStats), "guild_stats", len(guildStats), "player_stats", len(playerStats), "queues", len(queues))
}

func (p *BattleboardPoller) collectBattleSummaries(battles []tasks.Battle) ([]postgres.BattleSummary){
	summary := make([]postgres.BattleSummary, 0, len(battles))

	for _, battle := range battles {
		allianceSlice := make([]tasks.BattleAlliance, 0, len(battle.Alliances))
		for _, alliance := range battle.Alliances {
			allianceSlice = append(allianceSlice, alliance)
		}
		sort.Slice(allianceSlice, func(i, j int) bool {
			return allianceSlice[i].Kills > allianceSlice[j].Kills
		})
		allianceNames := make([]string, 0, len(battle.Alliances))
		for _, alliance := range allianceSlice {
			allianceNames = append(allianceNames, alliance.Name)
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
			guildNames = append(guildNames, guild.Name)
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