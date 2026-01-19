package battle_poller

import (
	"albionstats/internal/postgres"
	"albionstats/internal/tasks"
	"log/slog"
	"time"
)

type Config struct {
	APIClient       *tasks.Client
	Postgres        *postgres.Postgres
	Logger          *slog.Logger
	Region          string
}

type BattlePoller struct {
	apiClient *tasks.Client
	postgres *postgres.Postgres
	log *slog.Logger
	region string
}

func NewBattlePoller(cfg Config) (*BattlePoller) {
	return &BattlePoller{
		apiClient: cfg.APIClient,
		postgres: cfg.Postgres,
		log: cfg.Logger.With("component", "battle_poller", "region", cfg.Region),
		region: cfg.Region,
	}
}

func (p *BattlePoller) Run() {
	p.log.Info("battle polling started")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	p.runBatch() // Run once immediately

	for range ticker.C {
		p.runBatch()
	}
}

func (p *BattlePoller) runBatch() {
	queues, err := p.postgres.GetBattleQueuesByRegion(postgres.Region(p.region), 1)
	if err != nil {
		p.log.Error("get battle queues by region failed", "err", err)
		return
	}

	if len(queues) == 0 {
		p.log.Info("no battle queues found")
		time.Sleep(time.Second)
		return
	}

	for _, queue := range queues {
		events, err := p.fetchBattleEvents(queue.BattleID)
		if err != nil {
			p.log.Error("fetch battle events failed", "err", err)
			continue
		}

		allianceStats := p.processBattleAllianceStats(events)
		guildStats := p.processBattleGuildStats(events)
		playerStats := p.processPlayerStats(events)
		kills := p.processBattleKills(events)

		if err := p.postgres.UpdateBattleAllianceStats(allianceStats); err != nil {
			p.log.Error("update battle alliance stats failed", "err", err)
			continue
		}

		if err := p.postgres.UpdateBattleGuildStats(guildStats); err != nil {
			p.log.Error("insert battle guild stats failed", "err", err)
			continue
		}

		if err := p.postgres.UpdateBattlePlayerStats(playerStats); err != nil {
			p.log.Error("insert battle player stats failed", "err", err)
			continue
		}

		if err := p.postgres.InsertBattleKills(kills); err != nil {
			p.log.Error("insert battle kills failed", "err", err)
			continue
		}

		if err := p.postgres.DeleteBattleQueue(postgres.Region(p.region), queue.BattleID); err != nil {
			p.log.Error("delete battle queue failed", "err", err)
			continue
		}

		p.log.Info("battle processed", "battle_id", queue.BattleID, 
		    "alliance_stats", len(allianceStats), 
			"guild_stats", len(guildStats), 
			"player_stats", len(playerStats), 
			"kills", len(kills))
	}
}

func (p *BattlePoller) fetchBattleEvents(battleId int64) ([]tasks.Event, error) {
	var allEvents []tasks.Event
	offset := 0
	limit := 51

	for {
		events, err := p.apiClient.FetchBattleEvents(p.region, battleId, offset, limit)
		if err != nil {
			return nil, err
		}

		allEvents = append(allEvents, events...)

		// If we got fewer events than the limit, we've reached the end
		if len(events) < limit {
			break
		}

		// Otherwise, increment offset by 50 to get the next page
		offset += 50
	}

	return allEvents, nil
}

func (p *BattlePoller) processBattleAllianceStats(events []tasks.Event) []postgres.BattleAllianceStats {
	battleId := events[0].BattleID

	allianceTotalPlayers := make(map[string]int32)
	allianceTotalIp := make(map[string]float64)
	allianceDeathFame := make(map[string]int64)
	for _, event := range events {
		if event.Victim.AllianceName == "" {
			continue
		}
		allianceTotalPlayers[event.Victim.AllianceName] += 1
		allianceDeathFame[event.Victim.AllianceName] += event.TotalVictimKillFame
		allianceTotalIp[event.Victim.AllianceName] += event.Victim.AverageItemPower
	}

	for _, event := range events {
		if event.Killer.AllianceName == "" {
			continue
		}
		allianceTotalPlayers[event.Killer.AllianceName] += 1
		allianceTotalIp[event.Killer.AllianceName] += event.Killer.AverageItemPower
	}

	var playerStats []postgres.BattleAllianceStats
	for alliance, _ := range allianceTotalPlayers {
		deathFame := allianceDeathFame[alliance]
		averageIp := int32(allianceTotalIp[alliance] / float64(allianceTotalPlayers[alliance]))
		playerStats = append(playerStats, postgres.BattleAllianceStats{
			Region: postgres.Region(p.region),
			BattleID: battleId,
			AllianceName: alliance,
			DeathFame: &deathFame,
			IP: &averageIp,
		})
	}

	return playerStats
}

func (p *BattlePoller) processBattleGuildStats(events []tasks.Event) []postgres.BattleGuildStats {
	battleId := events[0].BattleID

	guildTotalPlayers := make(map[string]int32)
	guildTotalIp := make(map[string]float64)
	guildDeathFame := make(map[string]int64)
	for _, event := range events {
		if event.Victim.GuildName == "" {
			continue
		}
		guildTotalPlayers[event.Victim.GuildName] += 1
		guildDeathFame[event.Victim.GuildName] += event.TotalVictimKillFame
		guildTotalIp[event.Victim.GuildName] += event.Victim.AverageItemPower
	}

	for _, event := range events {
		if event.Killer.GuildName == "" {
			continue
		}
		guildTotalPlayers[event.Killer.GuildName] += 1
		guildTotalIp[event.Killer.GuildName] += event.Killer.AverageItemPower
	}

	var guildStats []postgres.BattleGuildStats
	for guild, _ := range guildTotalPlayers {
		deathFame := guildDeathFame[guild]
		averageIp := int32(guildTotalIp[guild] / float64(guildTotalPlayers[guild]))
		guildStats = append(guildStats, postgres.BattleGuildStats{
			Region: postgres.Region(p.region),
			BattleID: battleId,
			GuildName: guild,
			DeathFame: &deathFame,
			IP: &averageIp,
		})
	}

	return guildStats
}

func (p *BattlePoller) processPlayerStats(events []tasks.Event) []postgres.BattlePlayerStats {
	battleId := events[0].BattleID

	playerIp := make(map[string]float64)
	playerDeathFame := make(map[string]int64)
	playerWeapon := make(map[string]string)
	playerDamage := make(map[string]int64)
	playerHeal := make(map[string]int64)

	// Iterate kills first
	for _, event := range events {
		if _, ok := playerIp[event.Killer.Name]; !ok {
			playerIp[event.Killer.Name] = event.Victim.AverageItemPower
			if event.Victim.Equipment != nil {
				if mainHand, exists := event.Killer.Equipment["MainHand"]; mainHand != nil && exists {
					playerWeapon[event.Killer.Name] = mainHand.Type
				}
			}
		}
	}

	// ...deaths second
	for _, event := range events {
		playerDeathFame[event.Victim.Name] += event.TotalVictimKillFame

		if _, ok := playerIp[event.Victim.Name]; !ok {
			playerIp[event.Victim.Name] = event.Victim.AverageItemPower
			if event.Victim.Equipment != nil {
				if mainHand, exists := event.Victim.Equipment["MainHand"]; mainHand != nil && exists {
					playerWeapon[event.Victim.Name] = mainHand.Type
				}
			}
		}
	}

	// and then participants
	for _, event := range events {
		for _, participant := range event.Participants {
			playerDamage[participant.Name] += int64(participant.DamageDone)
			playerHeal[participant.Name] += int64(participant.SupportHealingDone)
		}
	}

	playerStats := make([]postgres.BattlePlayerStats, 0)

	for name, _ := range playerIp {
		deathFame := playerDeathFame[name]
		ip := int32(playerIp[name])
		weapon := playerWeapon[name]
		damage := playerDamage[name]
		heal := playerHeal[name]

		playerStats = append(playerStats, postgres.BattlePlayerStats{
			Region: postgres.Region(p.region),
			BattleID: battleId,
			PlayerName: name,
			DeathFame: &deathFame,
			IP: &ip,
			Weapon: &weapon,
			Damage: &damage,
			Heal: &heal,
		})
	}

	return playerStats
}

func (p *BattlePoller) processBattleKills(events []tasks.Event) []postgres.BattleKills {
	playerStats := make([]postgres.BattleKills, 0)
	for _, event := range events {
		killerWeapon := ""
		if event.Killer.Equipment != nil {
			if mainHand, exists := event.Killer.Equipment["MainHand"]; mainHand != nil && exists {
				killerWeapon = mainHand.Type
			}
		}

		victimWeapon := ""
		if event.Victim.Equipment != nil {
			if mainHand, exists := event.Victim.Equipment["MainHand"]; mainHand != nil && exists {
				victimWeapon = mainHand.Type
			}
		}

		playerStats = append(playerStats, postgres.BattleKills{
			Region: postgres.Region(p.region),
			BattleID: event.BattleID,
			TS: event.TimeStamp,
			KillerName: event.Killer.Name,
			KillerIP: int32(event.Killer.AverageItemPower),
			KillerWeapon: killerWeapon,
			VictimName: event.Victim.Name,
			VictimIP: int32(event.Victim.AverageItemPower),
			VictimWeapon: victimWeapon,
			Fame: event.TotalVictimKillFame,
		})
	}
	return playerStats
}
