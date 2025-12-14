package tasks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"albionstats/internal/api"
	"albionstats/internal/database"

	"errors"

	"gorm.io/gorm"
)

type PlayerPollerConfig struct {
	Region      string
	PageSize    int
	RatePerSec  int
	UserAgent   string
	Workers     int
	HTTPTimeout time.Duration
}

type PlayerPoller struct {
	apiClient *api.Client
	db        *gorm.DB
	cfg       PlayerPollerConfig
}

type processResult struct {
	playerPoll  database.PlayerPoll
	updatePoll  database.PlayerPoll
	statsLatest database.PlayerStatsLatest
	snapshot    database.PlayerStatsSnapshot
	delete      bool
	err         error
}

func (p *PlayerPoller) workerCount(n int) int {
	if p.cfg.Workers > 0 {
		if p.cfg.Workers < n {
			return p.cfg.Workers
		}
		return n
	}
	if p.cfg.RatePerSec > 0 {
		if p.cfg.RatePerSec < n {
			return p.cfg.RatePerSec
		}
		return n
	}
	if n < 1 {
		return 1
	}
	return n
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func NewPlayerPoller(db *gorm.DB, cfg PlayerPollerConfig) *PlayerPoller {
	apiClient := api.NewClient(cfg.Region, cfg.UserAgent, cfg.HTTPTimeout)
	return &PlayerPoller{
		apiClient: apiClient,
		db:        db,
		cfg:       cfg,
	}
}

func (p *PlayerPoller) fetchPlayersToPoll(ctx context.Context) ([]database.PlayerPoll, error) {
	var players []database.PlayerPoll
	now := time.Now().UTC()
	if err := p.db.WithContext(ctx).
		Where("region = ? AND next_poll_at <= ?", p.cfg.Region, now).
		Order("next_poll_at ASC").
		Limit(p.cfg.PageSize).
		Find(&players).Error; err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	return players, nil
}

func (p *PlayerPoller) handleIdleState(ctx context.Context) error {
	log.Printf("[%s] no players to poll", p.cfg.Region)
	idle := time.Second
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(idle):
		return nil
	}
}

func (p *PlayerPoller) setupWorkers(ctx context.Context, players []database.PlayerPoll) (*time.Ticker, chan database.PlayerPoll, chan processResult, *sync.WaitGroup) {
	rate := time.Second / time.Duration(p.cfg.RatePerSec)
	ticker := time.NewTicker(rate)

	workerCount := p.workerCount(len(players))
	jobs := make(chan database.PlayerPoll)
	results := make(chan processResult, len(players))

	var wg sync.WaitGroup

	log.Printf("[%s] batch size=%d rate=%d/s workers=%d", p.cfg.Region, len(players), p.cfg.RatePerSec, workerCount)

	// Start worker goroutines
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pl := range jobs {
				log.Printf("[%s] worker fetching player_id=%s", p.cfg.Region, pl.PlayerID)
				// shared rate limiter
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}

				results <- p.processPlayer(ctx, pl)
			}
		}()
	}

	// Start job sender goroutine
	go func() {
		for _, pl := range players {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- pl:
			}
		}
		close(jobs)
	}()

	return ticker, jobs, results, &wg
}

func (p *PlayerPoller) processResults(results <-chan processResult) ([]database.PlayerPoll, []database.PlayerStatsLatest, []database.PlayerStatsSnapshot, []database.PlayerPoll, []database.PlayerPoll) {
	var updatePolls []database.PlayerPoll
	var statsLatest []database.PlayerStatsLatest
	var snapshots []database.PlayerStatsSnapshot
	var deletes []database.PlayerPoll
	var failures []database.PlayerPoll

	for res := range results {
		if res.err != nil {
			log.Printf("[%s] player=%s err=%v", p.cfg.Region, res.playerPoll.PlayerID, res.err)
			nextErr := res.playerPoll.ErrorCount + 1
			backoff := failureBackoff(nextErr)
			failures = append(failures, database.PlayerPoll{
				Region:                res.playerPoll.Region,
				PlayerID:              res.playerPoll.PlayerID,
				NextPollAt:            time.Now().UTC().Add(backoff),
				ErrorCount:            nextErr,
				LastEncountered:       res.playerPoll.LastEncountered,
				KillboardLastActivity: res.playerPoll.KillboardLastActivity,
				OtherLastActivity:     res.playerPoll.OtherLastActivity,
				LastPollAt:            res.playerPoll.LastPollAt,
			})
			continue
		}
		if res.delete {
			deletes = append(deletes, res.playerPoll)
			continue
		}
		updatePolls = append(updatePolls, res.updatePoll)
		statsLatest = append(statsLatest, res.statsLatest)
		snapshots = append(snapshots, res.snapshot)
	}

	return updatePolls, statsLatest, snapshots, deletes, failures
}

func (p *PlayerPoller) applyDatabaseChanges(ctx context.Context, updatePolls []database.PlayerPoll, statsLatest []database.PlayerStatsLatest, snapshots []database.PlayerStatsSnapshot, deletes []database.PlayerPoll, failures []database.PlayerPoll) {
	if err := database.ApplyPlayerPollerDatabaseChanges(ctx, p.db, deletes, updatePolls, statsLatest, snapshots, failures); err != nil {
		log.Printf("[%s] database changes error: %v", p.cfg.Region, err)
	}
}

func (p *PlayerPoller) Run(ctx context.Context) {
	log.Printf("[%s] starting continuous polling", p.cfg.Region)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] stopped", p.cfg.Region)
			return
		default:
		}

		err := p.runBatch(ctx)
		if err != nil {
			log.Printf("[%s] batch error: %v", p.cfg.Region, err)
		}
	}
}

func (p *PlayerPoller) runBatch(ctx context.Context) error {
	players, err := p.fetchPlayersToPoll(ctx)
	if err != nil {
		return err
	}
	if len(players) == 0 {
		return p.handleIdleState(ctx)
	}

	ticker, _, results, wg := p.setupWorkers(ctx, players)
	defer ticker.Stop()

	wg.Wait()
	close(results)

	updatePolls, statsLatest, snapshots, deletes, failures := p.processResults(results)

	p.applyDatabaseChanges(ctx, updatePolls, statsLatest, snapshots, deletes, failures)

	return nil
}

func (p *PlayerPoller) processPlayer(ctx context.Context, pl database.PlayerPoll) processResult {
	ts := time.Now().UTC()
	resp, err := p.apiClient.FetchPlayer(ctx, pl.PlayerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return processResult{playerPoll: pl, delete: true}
		}
		return processResult{playerPoll: pl, err: fmt.Errorf("fetch player %s: %w", pl.PlayerID, err)}
	}

	apiTS := resp.LifetimeStatistics.Timestamp
	if apiTS == nil {
		return processResult{playerPoll: pl, delete: true}
	}
	nextPollAt, err := scheduleNextPoll(pl.LastEncountered, pl.KillboardLastActivity, pl.OtherLastActivity, ts)
	if err != nil {
		return processResult{playerPoll: pl, err: fmt.Errorf("schedule next poll: %w", err)}
	}

	// Create stats latest record
	statsLatest := database.PlayerStatsLatest{
		Region:                pl.Region,
		PlayerID:              pl.PlayerID,
		TS:                    ts,
		LastEncountered:       pl.LastEncountered,
		KillboardLastActivity: pl.KillboardLastActivity,
		OtherLastActivity:     resp.LifetimeStatistics.Timestamp,
		Name:                  resp.Name,
		GuildID:               api.NullableString(resp.GuildId),
		GuildName:             api.NullableString(resp.GuildName),
		AllianceID:            api.NullableString(resp.AllianceId),
		AllianceName:          api.NullableString(resp.AllianceName),
		AllianceTag:           api.NullableString(resp.AllianceTag),
		KillFame:              resp.KillFame,
		DeathFame:             resp.DeathFame,
		FameRatio:             api.NullableFloat64(resp.FameRatio),
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

	// Create snapshot record
	snapshot := database.PlayerStatsSnapshot{
		Region:                pl.Region,
		PlayerID:              pl.PlayerID,
		TS:                    ts,
		LastEncountered:       pl.LastEncountered,
		KillboardLastActivity: pl.KillboardLastActivity,
		OtherLastActivity:     resp.LifetimeStatistics.Timestamp,
		Name:                  resp.Name,
		GuildID:               api.NullableString(resp.GuildId),
		GuildName:             api.NullableString(resp.GuildName),
		AllianceID:            api.NullableString(resp.AllianceId),
		AllianceName:          api.NullableString(resp.AllianceName),
		AllianceTag:           api.NullableString(resp.AllianceTag),
		KillFame:              resp.KillFame,
		DeathFame:             resp.DeathFame,
		FameRatio:             api.NullableFloat64(resp.FameRatio),
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

	// Update player poll record
	updatePoll := database.PlayerPoll{
		Region:                pl.Region,
		PlayerID:              pl.PlayerID,
		LastPollAt:            &ts,
		NextPollAt:            nextPollAt,
		ErrorCount:            0,
		OtherLastActivity:     resp.LifetimeStatistics.Timestamp,
		LastEncountered:       pl.LastEncountered,
		KillboardLastActivity: pl.KillboardLastActivity,
	}

	return processResult{
		playerPoll:  pl,
		updatePoll:  updatePoll,
		statsLatest: statsLatest,
		snapshot:    snapshot,
	}
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
