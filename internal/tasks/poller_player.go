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
	APIBase     string
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
	playerState database.PlayerState
	updateState database.PlayerState
	snapshot    database.PlayerStatsSnapshot
	delete      bool
	err         error
	nextPollAt  time.Time
	priority    int
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
	apiClient := api.NewClient(cfg.APIBase, cfg.UserAgent, cfg.HTTPTimeout)
	return &PlayerPoller{
		apiClient: apiClient,
		db:        db,
		cfg:       cfg,
	}
}

func (p *PlayerPoller) fetchPlayersToPoll(ctx context.Context) ([]database.PlayerState, error) {
	var players []database.PlayerState
	now := time.Now().UTC()
	if err := p.db.WithContext(ctx).
		Where("next_poll_at <= ?", now).
		Order("priority ASC, next_poll_at ASC").
		Limit(p.cfg.PageSize).
		Find(&players).Error; err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	return players, nil
}

func (p *PlayerPoller) handleIdleState(ctx context.Context) error {
	log.Printf("no players to poll")
	idle := time.Second
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(idle):
		return nil
	}
}

func (p *PlayerPoller) setupWorkers(ctx context.Context, players []database.PlayerState) (*time.Ticker, chan database.PlayerState, chan processResult, *sync.WaitGroup) {
	rate := time.Second / time.Duration(p.cfg.RatePerSec)
	ticker := time.NewTicker(rate)

	workerCount := p.workerCount(len(players))
	jobs := make(chan database.PlayerState)
	results := make(chan processResult, len(players))

	var wg sync.WaitGroup

	log.Printf("batch size=%d rate=%d/s workers=%d", len(players), p.cfg.RatePerSec, workerCount)

	// Start worker goroutines
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pl := range jobs {
				log.Printf("worker fetching player_id=%s name=%s", pl.PlayerID, pl.Name)
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

func (p *PlayerPoller) processResults(results <-chan processResult) ([]database.PlayerState, []database.PlayerStatsSnapshot, []database.PlayerState, []database.PlayerState) {
	var updates []database.PlayerState
	var snapshots []database.PlayerStatsSnapshot
	var deletes []database.PlayerState
	var failures []database.PlayerState

	for res := range results {
		if res.err != nil {
			log.Printf("player=%s err=%v", res.playerState.PlayerID, res.err)
			nextErr := res.playerState.ErrorCount + 1
			backoff := failureBackoff(nextErr)
			failures = append(failures, database.PlayerState{
				Region:     res.playerState.Region,
				PlayerID:   res.playerState.PlayerID,
				ErrorCount: nextErr,
				NextPollAt: time.Now().UTC().Add(backoff),
			})
			continue
		}
		if res.delete {
			deletes = append(deletes, res.playerState)
			continue
		}
		updates = append(updates, res.updateState)
		snapshots = append(snapshots, res.snapshot)
	}

	return updates, snapshots, deletes, failures
}

func (p *PlayerPoller) applyDatabaseChanges(ctx context.Context, updates []database.PlayerState, snapshots []database.PlayerStatsSnapshot, deletes []database.PlayerState, failures []database.PlayerState) {
	// apply deletes (grouped by region)
	if len(deletes) > 0 {
		log.Printf("deleting %d players", len(deletes))
		regionBuckets := make(map[string][]string)
		for _, d := range deletes {
			regionBuckets[d.Region] = append(regionBuckets[d.Region], d.PlayerID)
		}
		for region, ids := range regionBuckets {
			if err := p.db.WithContext(ctx).Delete(&database.PlayerState{}, "region = ? AND player_id IN ?", region, ids).Error; err != nil {
				log.Printf("delete err region=%s ids=%d: %v", region, len(ids), err)
			}
		}
	}

	// upsert player_state
	if len(updates) > 0 {
		log.Printf("upserting %d players", len(updates))
		if err := database.BulkUpsertStates(ctx, p.db, updates); err != nil {
			log.Printf("upsert states err: %v", err)
		}
	}

	// apply failures (increment error_count, backoff)
	if len(failures) > 0 {
		log.Printf("upserting %d errors", len(failures))
		if err := database.BulkUpsertFailures(ctx, p.db, failures); err != nil {
			log.Printf("failure upsert err: %v", err)
		}
	}

	// insert snapshots
	if len(snapshots) > 0 {
		log.Printf("inserting %d snapshots", len(snapshots))
		if err := p.db.WithContext(ctx).Create(&snapshots).Error; err != nil {
			log.Printf("insert snapshots err: %v", err)
		}
	}
}

func (p *PlayerPoller) Run(ctx context.Context) {
	log.Printf("starting continuous polling")

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopped")
			return
		default:
		}

		err := p.runBatch(ctx)
		if err != nil {
			log.Printf("batch error: %v", err)
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

	updates, snapshots, deletes, failures := p.processResults(results)

	p.applyDatabaseChanges(ctx, updates, snapshots, deletes, failures)

	return nil
}

func (p *PlayerPoller) processPlayer(ctx context.Context, pl database.PlayerState) processResult {
	ts := time.Now().UTC()
	resp, err := p.apiClient.FetchPlayer(ctx, pl.PlayerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return processResult{playerState: pl, delete: true}
		}
		return processResult{playerState: pl, err: fmt.Errorf("fetch player %s: %w", pl.PlayerID, err)}
	}

	apiTS := resp.LifetimeStatistics.Timestamp
	if apiTS == nil {
		return processResult{playerState: pl, delete: true}
	}
	nextPollAt, priority := scheduleNextPoll(*apiTS, ts)

	snapshot := database.PlayerStatsSnapshot{
		Region:              pl.Region,
		PlayerID:            pl.PlayerID,
		TS:                  ts,
		APITimestamp:        apiTS,
		Name:                resp.Name,
		GuildID:             api.NullableString(resp.GuildId),
		GuildName:           api.NullableString(resp.GuildName),
		AllianceID:          api.NullableString(resp.AllianceId),
		AllianceName:        api.NullableString(resp.AllianceName),
		AllianceTag:         api.NullableString(resp.AllianceTag),
		KillFame:            api.NullableInt64(resp.KillFame),
		DeathFame:           api.NullableInt64(resp.DeathFame),
		FameRatio:           api.NullableFloat64(resp.FameRatio),
		CraftingTotal:       api.NullableInt64(resp.LifetimeStatistics.Crafting.Total),
		CraftingRoyal:       api.NullableInt64(resp.LifetimeStatistics.Crafting.Royal),
		CraftingOutlands:    api.NullableInt64(resp.LifetimeStatistics.Crafting.Outlands),
		CraftingAvalon:      api.NullableInt64(resp.LifetimeStatistics.Crafting.Avalon),
		FishingFame:         api.NullableInt64(resp.FishingFame),
		FarmingFame:         api.NullableInt64(resp.FarmingFame),
		CrystalLeagueFame:   api.NullableInt64(resp.CrystalLeague),
		PveTotal:            api.NullableInt64(resp.LifetimeStatistics.PvE.Total),
		PveRoyal:            api.NullableInt64(resp.LifetimeStatistics.PvE.Royal),
		PveOutlands:         api.NullableInt64(resp.LifetimeStatistics.PvE.Outlands),
		PveAvalon:           api.NullableInt64(resp.LifetimeStatistics.PvE.Avalon),
		PveHellgate:         api.NullableInt64(resp.LifetimeStatistics.PvE.Hellgate),
		PveCorrupted:        api.NullableInt64(resp.LifetimeStatistics.PvE.CorruptedDungeon),
		PveMists:            api.NullableInt64(resp.LifetimeStatistics.PvE.Mists),
		GatherFiberTotal:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Total),
		GatherFiberRoyal:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Royal),
		GatherFiberOutlands: api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Outlands),
		GatherFiberAvalon:   api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Avalon),
		GatherHideTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Total),
		GatherHideRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Royal),
		GatherHideOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Outlands),
		GatherHideAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Avalon),
		GatherOreTotal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Total),
		GatherOreRoyal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Royal),
		GatherOreOutlands:   api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Outlands),
		GatherOreAvalon:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Avalon),
		GatherRockTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Total),
		GatherRockRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Royal),
		GatherRockOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Outlands),
		GatherRockAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Avalon),
		GatherWoodTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Total),
		GatherWoodRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Royal),
		GatherWoodOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Outlands),
		GatherWoodAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Avalon),
		GatherAllTotal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Total),
		GatherAllRoyal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Royal),
		GatherAllOutlands:   api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Outlands),
		GatherAllAvalon:     api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Avalon),
	}

	// Update player_state identity and activity tracking
	update := database.PlayerState{
		Region:              pl.Region,
		PlayerID:            pl.PlayerID,
		Name:                resp.Name,
		GuildID:             api.NullableString(resp.GuildId),
		GuildName:           api.NullableString(resp.GuildName),
		AllianceID:          api.NullableString(resp.AllianceId),
		AllianceName:        api.NullableString(resp.AllianceName),
		AllianceTag:         api.NullableString(resp.AllianceTag),
		LastPolled:          &ts,
		LastSeen:            apiTS,
		NextPollAt:          nextPollAt,
		Priority:            priority,
		ErrorCount:          0,
		LastError:           nil,
		KillFame:            api.NullableInt64(resp.KillFame),
		DeathFame:           api.NullableInt64(resp.DeathFame),
		FameRatio:           api.NullableFloat64(resp.FameRatio),
		PveTotal:            api.NullableInt64(resp.LifetimeStatistics.PvE.Total),
		PveRoyal:            api.NullableInt64(resp.LifetimeStatistics.PvE.Royal),
		PveOutlands:         api.NullableInt64(resp.LifetimeStatistics.PvE.Outlands),
		PveAvalon:           api.NullableInt64(resp.LifetimeStatistics.PvE.Avalon),
		PveHellgate:         api.NullableInt64(resp.LifetimeStatistics.PvE.Hellgate),
		PveCorrupted:        api.NullableInt64(resp.LifetimeStatistics.PvE.CorruptedDungeon),
		PveMists:            api.NullableInt64(resp.LifetimeStatistics.PvE.Mists),
		GatherFiberTotal:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Total),
		GatherFiberRoyal:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Royal),
		GatherFiberOutlands: api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Outlands),
		GatherFiberAvalon:   api.NullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Avalon),
		GatherHideTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Total),
		GatherHideRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Royal),
		GatherHideOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Outlands),
		GatherHideAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Hide.Avalon),
		GatherOreTotal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Total),
		GatherOreRoyal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Royal),
		GatherOreOutlands:   api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Outlands),
		GatherOreAvalon:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Ore.Avalon),
		GatherRockTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Total),
		GatherRockRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Royal),
		GatherRockOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Outlands),
		GatherRockAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Rock.Avalon),
		GatherWoodTotal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Total),
		GatherWoodRoyal:     api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Royal),
		GatherWoodOutlands:  api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Outlands),
		GatherWoodAvalon:    api.NullableInt64(resp.LifetimeStatistics.Gathering.Wood.Avalon),
		GatherAllTotal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Total),
		GatherAllRoyal:      api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Royal),
		GatherAllOutlands:   api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Outlands),
		GatherAllAvalon:     api.NullableInt64(resp.LifetimeStatistics.Gathering.All.Avalon),
		CraftingTotal:       api.NullableInt64(resp.LifetimeStatistics.Crafting.Total),
		CraftingRoyal:       api.NullableInt64(resp.LifetimeStatistics.Crafting.Royal),
		CraftingOutlands:    api.NullableInt64(resp.LifetimeStatistics.Crafting.Outlands),
		CraftingAvalon:      api.NullableInt64(resp.LifetimeStatistics.Crafting.Avalon),
		FishingFame:         api.NullableInt64(resp.FishingFame),
		FarmingFame:         api.NullableInt64(resp.FarmingFame),
		CrystalLeagueFame:   api.NullableInt64(resp.CrystalLeague),
	}

	return processResult{
		playerState: pl,
		updateState: update,
		snapshot:    snapshot,
		priority:    priority,
		nextPollAt:  nextPollAt,
	}
}

func scheduleNextPoll(apiTS, now time.Time) (time.Time, int) {
	staleness := now.Sub(apiTS)
	switch {
	case staleness <= 24*time.Hour:
		return now.Add(6 * time.Hour), 200
	case staleness <= 7*24*time.Hour:
		return now.Add(24 * time.Hour), 300
	case staleness <= 30*24*time.Hour:
		return now.Add(48 * time.Hour), 400
	default:
		return now.Add(24 * 30 * time.Hour), 500
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
