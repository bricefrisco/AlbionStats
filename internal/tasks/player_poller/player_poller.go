package player_poller

import (
	"albionstats/internal/api"
	"albionstats/internal/sqlite"
	"albionstats/internal/util"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	APIClient  *api.Client
	SQLite     *sqlite.SQLite
	Logger     *slog.Logger
	Region     string
	BatchSize  int
	RatePerSec int
}

type PlayerPoller struct {
	api        *api.Client
	sqlite     *sqlite.SQLite
	log        *slog.Logger
	region     string
	batchSize  int
	ratePerSec int
}

type processResult struct {
	poll             sqlite.PlayerPoll
	stats            sqlite.PlayerStats
	shouldDeletePoll bool
}

func NewPlayerPoller(cfg Config) (*PlayerPoller, error) {
	return &PlayerPoller{
		api:        cfg.APIClient,
		sqlite:     cfg.SQLite,
		log:        cfg.Logger.With("component", "player_poller", "region", cfg.Region),
		region:     cfg.Region,
		batchSize:  cfg.BatchSize,
		ratePerSec: cfg.RatePerSec,
	}, nil
}

func (p *PlayerPoller) Run() {
	p.log.Info("player polling started", "region", p.region, "batch_size", p.batchSize, "rate_per_sec", p.ratePerSec)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		p.runBatch()
	}
}

func (p *PlayerPoller) runBatch() {
	players, err := p.sqlite.FetchPlayersToPoll(p.region, p.batchSize)
	if err != nil {
		p.log.Error("fetch players to poll failed", "err", err)
		return
	}

	if len(players) == 0 {
		time.Sleep(time.Second)
	}

	p.log.Info("starting batch", "num_players", len(players))
	pool := NewWorkerPool[processResult](p.ratePerSec)
	for _, player := range players {
		pool.Add(func() processResult {
			return p.processPlayer(player)
		})
	}

	results := pool.ExecuteJobs()
	p.processResults(results)
}

func (p *PlayerPoller) processPlayer(player sqlite.PlayerPoll) processResult {
	now := time.Now().UTC()

	resp, err := p.api.FetchPlayer(p.region, player.PlayerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return processResult{shouldDeletePoll: true}
		}

		p.log.Warn("player poll failed", "player_id", player.PlayerID, "err", err.Error())
		return processResult{poll: sqlite.PlayerPoll{
			Region:                player.Region,
			PlayerID:              player.PlayerID,
			LastPollAt:            player.LastPollAt,
			NextPollAt:            time.Now().UTC().Add(failureBackoff(player.ErrorCount + 1)),
			ErrorCount:            player.ErrorCount + 1,
			LastEncountered:       player.LastEncountered,
			KillboardLastActivity: player.KillboardLastActivity,
			OtherLastActivity:     player.OtherLastActivity,
		}}
	}

	if resp.LifetimeStatistics.Timestamp == nil {
		return processResult{shouldDeletePoll: true}
	}

	nextPollAt, err := scheduleNextPoll(player.LastEncountered, player.KillboardLastActivity, player.OtherLastActivity, now)
	if err != nil {
		p.log.Warn("player poll failed", "player_id", player.PlayerID, "err", err.Error())
		return processResult{poll: sqlite.PlayerPoll{
			Region:                player.Region,
			PlayerID:              player.PlayerID,
			LastPollAt:            player.LastPollAt,
			NextPollAt:            time.Now().UTC().Add(failureBackoff(player.ErrorCount + 1)),
			ErrorCount:            player.ErrorCount + 1,
			LastEncountered:       player.LastEncountered,
			KillboardLastActivity: player.KillboardLastActivity,
			OtherLastActivity:     player.OtherLastActivity,
		}}
	}

	poll := sqlite.PlayerPoll{
		Region:                player.Region,
		PlayerID:              player.PlayerID,
		LastPollAt:            &now,
		NextPollAt:            nextPollAt,
		ErrorCount:            0,
		OtherLastActivity:     resp.LifetimeStatistics.Timestamp,
		LastEncountered:       player.LastEncountered,
		KillboardLastActivity: player.KillboardLastActivity,
	}

	stats := sqlite.PlayerStats{
		Region:                player.Region,
		PlayerID:              player.PlayerID,
		TS:                    now,
		LastEncountered:       player.LastEncountered,
		KillboardLastActivity: player.KillboardLastActivity,
		OtherLastActivity:     resp.LifetimeStatistics.Timestamp,
		Name:                  resp.Name,
		GuildID:               util.NullableString(resp.GuildId),
		GuildName:             util.NullableString(resp.GuildName),
		AllianceID:            util.NullableString(resp.AllianceId),
		AllianceName:          util.NullableString(resp.AllianceName),
		AllianceTag:           util.NullableString(resp.AllianceTag),
		KillFame:              resp.KillFame,
		DeathFame:             resp.DeathFame,
		FameRatio:             util.NullableFloat64(resp.FameRatio),
		PveTotal:              resp.LifetimeStatistics.PvE.Total,
		PveRoyal:              resp.LifetimeStatistics.PvE.Royal,
		PveOutlands:           resp.LifetimeStatistics.PvE.Outlands,
		PveAvalon:             resp.LifetimeStatistics.PvE.Avalon,
		PveHellgate:           resp.LifetimeStatistics.PvE.Hellgate,
		PveCorrupted:          resp.LifetimeStatistics.PvE.CorruptedDungeon,
		PveMists:              resp.LifetimeStatistics.PvE.Mists,
		GatherFiberTotal:      resp.LifetimeStatistics.Gathering.Fiber.Total,
		GatherFiberRoyal:      resp.LifetimeStatistics.Gathering.Fiber.Royal,
		GatherFiberOutlands:   resp.LifetimeStatistics.Gathering.Fiber.Outlands,
		GatherFiberAvalon:     resp.LifetimeStatistics.Gathering.Fiber.Avalon,
		GatherHideTotal:       resp.LifetimeStatistics.Gathering.Hide.Total,
		GatherHideRoyal:       resp.LifetimeStatistics.Gathering.Hide.Royal,
		GatherHideOutlands:    resp.LifetimeStatistics.Gathering.Hide.Outlands,
		GatherHideAvalon:      resp.LifetimeStatistics.Gathering.Hide.Avalon,
		GatherOreTotal:        resp.LifetimeStatistics.Gathering.Ore.Total,
		GatherOreRoyal:        resp.LifetimeStatistics.Gathering.Ore.Royal,
		GatherOreOutlands:     resp.LifetimeStatistics.Gathering.Ore.Outlands,
		GatherOreAvalon:       resp.LifetimeStatistics.Gathering.Ore.Avalon,
		GatherRockTotal:       resp.LifetimeStatistics.Gathering.Rock.Total,
		GatherRockRoyal:       resp.LifetimeStatistics.Gathering.Rock.Royal,
		GatherRockOutlands:    resp.LifetimeStatistics.Gathering.Rock.Outlands,
		GatherRockAvalon:      resp.LifetimeStatistics.Gathering.Rock.Avalon,
		GatherWoodTotal:       resp.LifetimeStatistics.Gathering.Wood.Total,
		GatherWoodRoyal:       resp.LifetimeStatistics.Gathering.Wood.Royal,
		GatherWoodOutlands:    resp.LifetimeStatistics.Gathering.Wood.Outlands,
		GatherWoodAvalon:      resp.LifetimeStatistics.Gathering.Wood.Avalon,
		GatherAllTotal:        resp.LifetimeStatistics.Gathering.All.Total,
		GatherAllRoyal:        resp.LifetimeStatistics.Gathering.All.Royal,
		GatherAllOutlands:     resp.LifetimeStatistics.Gathering.All.Outlands,
		GatherAllAvalon:       resp.LifetimeStatistics.Gathering.All.Avalon,
		CraftingTotal:         resp.LifetimeStatistics.Crafting.Total,
		CraftingRoyal:         resp.LifetimeStatistics.Crafting.Royal,
		CraftingOutlands:      resp.LifetimeStatistics.Crafting.Outlands,
		CraftingAvalon:        resp.LifetimeStatistics.Crafting.Avalon,
		FishingFame:           resp.FishingFame,
		FarmingFame:           resp.FarmingFame,
		CrystalLeagueFame:     resp.CrystalLeague,
	}

	return processResult{poll: poll, stats: stats}
}

func (p *PlayerPoller) processResults(results []processResult) {
	deletes := make([]sqlite.PlayerPoll, 0)
	polls := make([]sqlite.PlayerPoll, 0)
	stats := make([]sqlite.PlayerStats, 0)
	for _, result := range results {
		if result.shouldDeletePoll {
			deletes = append(deletes, result.poll)
		} else {
			polls = append(polls, result.poll)
			stats = append(stats, result.stats)
		}
	}

	if err := p.sqlite.DeletePlayerPolls(deletes); err != nil {
		p.log.Error("delete player polls failed", "err", err)
		return
	}

	if err := p.sqlite.UpdatePlayerPolls(polls); err != nil {
		p.log.Error("update player polls failed", "err", err)
		return
	}

	if err := p.sqlite.UpsertPlayerStats(stats); err != nil {
		p.log.Error("upsert player stats failed", "err", err)
		return
	}

	p.log.Info("processed results", "num_deletes", len(deletes), "num_polls", len(polls), "num_stats", len(stats))
}

func scheduleNextPoll(lastEncountered, killboardLastActivity, otherLastActivity *time.Time, now time.Time) (time.Time, error) {
	// Find the most recent activity timestamp among the three
	var mostRecent *time.Time
	if lastEncountered != nil {
		mostRecent = lastEncountered
	}
	if killboardLastActivity != nil && (mostRecent == nil || killboardLastActivity.After(*mostRecent)) {
		mostRecent = killboardLastActivity
	}
	if otherLastActivity != nil && (mostRecent == nil || otherLastActivity.After(*mostRecent)) {
		mostRecent = otherLastActivity
	}

	// If no activity timestamps are available, this should never happen
	if mostRecent == nil {
		return time.Time{}, fmt.Errorf("no activity timestamps available for player")
	}

	staleness := now.Sub(*mostRecent)
	switch {
	case staleness <= 24*time.Hour:
		return now.Add(6 * time.Hour), nil
	case staleness <= 7*24*time.Hour:
		return now.Add(24 * time.Hour), nil
	case staleness <= 30*24*time.Hour:
		return now.Add(48 * time.Hour), nil
	default:
		return now.Add(24 * 30 * time.Hour), nil
	}
}

func failureBackoff(errorCount int) time.Duration {
	base := 15 * time.Second
	maxBackoff := 24 * time.Hour
	// exponential backoff capped
	shift := errorCount
	if shift > 6 {
		shift = 6
	}
	backoff := base * (1 << shift)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}
