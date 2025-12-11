package playerpoller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"albionstats/internal/models"

	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	APIBase    string
	PageSize   int
	RatePerSec int
	UserAgent  string
	Workers    int
}

type Poller struct {
	client *http.Client
	db     *gorm.DB
	cfg    Config
}

type processResult struct {
	playerState models.PlayerState
	updateState models.PlayerState
	snapshot    models.PlayerStatsSnapshot
	delete      bool
	err         error
	nextPollAt  time.Time
	priority    int
}

func (p *Poller) workerCount(n int) int {
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

func New(client *http.Client, db *gorm.DB, cfg Config) *Poller {
	return &Poller{client: client, db: db, cfg: cfg}
}

func (p *Poller) Run(ctx context.Context) error {
	var players []models.PlayerState
	now := time.Now().UTC()
	if err := p.db.WithContext(ctx).
		Where("next_poll_at <= ?", now).
		Order("priority ASC, next_poll_at ASC").
		Limit(p.cfg.PageSize).
		Find(&players).Error; err != nil {
		return fmt.Errorf("query players: %w", err)
	}
	if len(players) == 0 {
		log.Printf("player-poller: no players to poll")
		idle := time.Second
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(idle):
			return nil
		}
	}

	rate := time.Second / time.Duration(p.cfg.RatePerSec)
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	workerCount := p.workerCount(len(players))
	jobs := make(chan models.PlayerState)
	results := make(chan processResult, len(players))

	var wg sync.WaitGroup

	log.Printf("player-poller: batch size=%d rate=%d/s workers=%d", len(players), p.cfg.RatePerSec, workerCount)

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pl := range jobs {
				log.Printf("player-poller: worker fetching player_id=%s name=%s", pl.PlayerID, pl.Name)
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

	wg.Wait()
	close(results)

	var updates []models.PlayerState
	var snapshots []models.PlayerStatsSnapshot
	var deletes []models.PlayerState
	var failures []models.PlayerState

	for res := range results {
		if res.err != nil {
			log.Printf("player-poller: player=%s err=%v", res.playerState.PlayerID, res.err)
			nextErr := res.playerState.ErrorCount + 1
			backoff := failureBackoff(nextErr)
			failures = append(failures, models.PlayerState{
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

	// apply deletes (grouped by region)
	if len(deletes) > 0 {
		log.Printf("player-poller: deleting %d players", len(deletes))
		regionBuckets := make(map[string][]string)
		for _, d := range deletes {
			regionBuckets[d.Region] = append(regionBuckets[d.Region], d.PlayerID)
		}
		for region, ids := range regionBuckets {
			if err := p.db.WithContext(ctx).Delete(&models.PlayerState{}, "region = ? AND player_id IN ?", region, ids).Error; err != nil {
				log.Printf("player-poller: delete err region=%s ids=%d: %v", region, len(ids), err)
			}
		}
	}

	// upsert player_state
	if len(updates) > 0 {
		log.Printf("player-poller: upserting %d players", len(updates))
		if err := p.bulkUpsertStates(ctx, updates); err != nil {
			log.Printf("player-poller: upsert states err: %v", err)
		}
	}

	// apply failures (increment error_count, backoff)
	if len(failures) > 0 {
		log.Printf("player-poller: upserting %d errors", len(failures))
		if err := p.bulkUpsertFailures(ctx, failures); err != nil {
			log.Printf("player-poller: failure upsert err: %v", err)
		}
	}

	// insert snapshots
	if len(snapshots) > 0 {
		log.Printf("player-poller: inserting %d snapshots", len(snapshots))
		if err := p.db.WithContext(ctx).Create(&snapshots).Error; err != nil {
			log.Printf("player-poller: insert snapshots err: %v", err)
		}
	}

	return nil
}

func (p *Poller) processPlayer(ctx context.Context, pl models.PlayerState) processResult {
	ts := time.Now().UTC()
	resp, err := p.fetchPlayer(ctx, pl.PlayerID)
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

	snapshot := models.PlayerStatsSnapshot{
		Region:              pl.Region,
		PlayerID:            pl.PlayerID,
		TS:                  ts,
		APITimestamp:        apiTS,
		Name:                resp.Name,
		GuildID:             nullableString(resp.GuildId),
		GuildName:           nullableString(resp.GuildName),
		AllianceID:          nullableString(resp.AllianceId),
		AllianceName:        nullableString(resp.AllianceName),
		AllianceTag:         nullableString(resp.AllianceTag),
		KillFame:            nullableInt64(resp.KillFame),
		DeathFame:           nullableInt64(resp.DeathFame),
		FameRatio:           nullableFloat64(resp.FameRatio),
		CraftingTotal:       nullableInt64(resp.LifetimeStatistics.Crafting.Total),
		CraftingRoyal:       nullableInt64(resp.LifetimeStatistics.Crafting.Royal),
		CraftingOutlands:    nullableInt64(resp.LifetimeStatistics.Crafting.Outlands),
		CraftingAvalon:      nullableInt64(resp.LifetimeStatistics.Crafting.Avalon),
		FishingFame:         nullableInt64(resp.FishingFame),
		FarmingFame:         nullableInt64(resp.FarmingFame),
		CrystalLeagueFame:   nullableInt64(resp.CrystalLeague),
		PveTotal:            nullableInt64(resp.LifetimeStatistics.PvE.Total),
		PveRoyal:            nullableInt64(resp.LifetimeStatistics.PvE.Royal),
		PveOutlands:         nullableInt64(resp.LifetimeStatistics.PvE.Outlands),
		PveAvalon:           nullableInt64(resp.LifetimeStatistics.PvE.Avalon),
		PveHellgate:         nullableInt64(resp.LifetimeStatistics.PvE.Hellgate),
		PveCorrupted:        nullableInt64(resp.LifetimeStatistics.PvE.CorruptedDungeon),
		PveMists:            nullableInt64(resp.LifetimeStatistics.PvE.Mists),
		GatherFiberTotal:    nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Total),
		GatherFiberRoyal:    nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Royal),
		GatherFiberOutlands: nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Outlands),
		GatherFiberAvalon:   nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Avalon),
		GatherHideTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Total),
		GatherHideRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Royal),
		GatherHideOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Outlands),
		GatherHideAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Avalon),
		GatherOreTotal:      nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Total),
		GatherOreRoyal:      nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Royal),
		GatherOreOutlands:   nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Outlands),
		GatherOreAvalon:     nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Avalon),
		GatherRockTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Total),
		GatherRockRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Royal),
		GatherRockOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Outlands),
		GatherRockAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Avalon),
		GatherWoodTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Total),
		GatherWoodRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Royal),
		GatherWoodOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Outlands),
		GatherWoodAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Avalon),
		GatherAllTotal:      nullableInt64(resp.LifetimeStatistics.Gathering.All.Total),
		GatherAllRoyal:      nullableInt64(resp.LifetimeStatistics.Gathering.All.Royal),
		GatherAllOutlands:   nullableInt64(resp.LifetimeStatistics.Gathering.All.Outlands),
		GatherAllAvalon:     nullableInt64(resp.LifetimeStatistics.Gathering.All.Avalon),
	}

	// Update player_state identity and activity tracking
	update := models.PlayerState{
		Region:              pl.Region,
		PlayerID:            pl.PlayerID,
		Name:                resp.Name,
		GuildID:             nullableString(resp.GuildId),
		GuildName:           nullableString(resp.GuildName),
		AllianceID:          nullableString(resp.AllianceId),
		AllianceName:        nullableString(resp.AllianceName),
		AllianceTag:         nullableString(resp.AllianceTag),
		LastPolled:          &ts,
		LastSeen:            apiTS,
		NextPollAt:          nextPollAt,
		Priority:            priority,
		ErrorCount:          0,
		LastError:           nil,
		KillFame:            nullableInt64(resp.KillFame),
		DeathFame:           nullableInt64(resp.DeathFame),
		FameRatio:           nullableFloat64(resp.FameRatio),
		PveTotal:            nullableInt64(resp.LifetimeStatistics.PvE.Total),
		PveRoyal:            nullableInt64(resp.LifetimeStatistics.PvE.Royal),
		PveOutlands:         nullableInt64(resp.LifetimeStatistics.PvE.Outlands),
		PveAvalon:           nullableInt64(resp.LifetimeStatistics.PvE.Avalon),
		PveHellgate:         nullableInt64(resp.LifetimeStatistics.PvE.Hellgate),
		PveCorrupted:        nullableInt64(resp.LifetimeStatistics.PvE.CorruptedDungeon),
		PveMists:            nullableInt64(resp.LifetimeStatistics.PvE.Mists),
		GatherFiberTotal:    nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Total),
		GatherFiberRoyal:    nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Royal),
		GatherFiberOutlands: nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Outlands),
		GatherFiberAvalon:   nullableInt64(resp.LifetimeStatistics.Gathering.Fiber.Avalon),
		GatherHideTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Total),
		GatherHideRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Royal),
		GatherHideOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Outlands),
		GatherHideAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Hide.Avalon),
		GatherOreTotal:      nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Total),
		GatherOreRoyal:      nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Royal),
		GatherOreOutlands:   nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Outlands),
		GatherOreAvalon:     nullableInt64(resp.LifetimeStatistics.Gathering.Ore.Avalon),
		GatherRockTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Total),
		GatherRockRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Royal),
		GatherRockOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Outlands),
		GatherRockAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Rock.Avalon),
		GatherWoodTotal:     nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Total),
		GatherWoodRoyal:     nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Royal),
		GatherWoodOutlands:  nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Outlands),
		GatherWoodAvalon:    nullableInt64(resp.LifetimeStatistics.Gathering.Wood.Avalon),
		GatherAllTotal:      nullableInt64(resp.LifetimeStatistics.Gathering.All.Total),
		GatherAllRoyal:      nullableInt64(resp.LifetimeStatistics.Gathering.All.Royal),
		GatherAllOutlands:   nullableInt64(resp.LifetimeStatistics.Gathering.All.Outlands),
		GatherAllAvalon:     nullableInt64(resp.LifetimeStatistics.Gathering.All.Avalon),
		CraftingTotal:       nullableInt64(resp.LifetimeStatistics.Crafting.Total),
		CraftingRoyal:       nullableInt64(resp.LifetimeStatistics.Crafting.Royal),
		CraftingOutlands:    nullableInt64(resp.LifetimeStatistics.Crafting.Outlands),
		CraftingAvalon:      nullableInt64(resp.LifetimeStatistics.Crafting.Avalon),
		FishingFame:         nullableInt64(resp.FishingFame),
		FarmingFame:         nullableInt64(resp.FarmingFame),
		CrystalLeagueFame:   nullableInt64(resp.CrystalLeague),
	}

	return processResult{
		playerState: pl,
		updateState: update,
		snapshot:    snapshot,
		priority:    priority,
		nextPollAt:  nextPollAt,
	}
}

func (p *Poller) bulkUpsertStates(ctx context.Context, states []models.PlayerState) error {
	if len(states) == 0 {
		return nil
	}

	always := map[string]interface{}{
		"name":          gorm.Expr("excluded.name"),
		"guild_id":      gorm.Expr("excluded.guild_id"),
		"guild_name":    gorm.Expr("excluded.guild_name"),
		"alliance_id":   gorm.Expr("excluded.alliance_id"),
		"alliance_name": gorm.Expr("excluded.alliance_name"),
		"alliance_tag":  gorm.Expr("excluded.alliance_tag"),
		"kill_fame":     gorm.Expr("excluded.kill_fame"),
		"death_fame":    gorm.Expr("excluded.death_fame"),
		"fame_ratio":    gorm.Expr("excluded.fame_ratio"),
		"last_polled":   gorm.Expr("excluded.last_polled"),
		"last_seen":     gorm.Expr("excluded.last_seen"),
		"next_poll_at":  gorm.Expr("excluded.next_poll_at"),
		"priority":      gorm.Expr("excluded.priority"),
		"error_count":   gorm.Expr("excluded.error_count"),
		"last_error":    gorm.Expr("excluded.last_error"),
	}

	statCols := []string{
		"pve_total", "pve_royal", "pve_outlands", "pve_avalon",
		"pve_hellgate", "pve_corrupted", "pve_mists",
		"gather_fiber_total", "gather_fiber_royal", "gather_fiber_outlands", "gather_fiber_avalon",
		"gather_hide_total", "gather_hide_royal", "gather_hide_outlands", "gather_hide_avalon",
		"gather_ore_total", "gather_ore_royal", "gather_ore_outlands", "gather_ore_avalon",
		"gather_rock_total", "gather_rock_royal", "gather_rock_outlands", "gather_rock_avalon",
		"gather_wood_total", "gather_wood_royal", "gather_wood_outlands", "gather_wood_avalon",
		"gather_all_total", "gather_all_royal", "gather_all_outlands", "gather_all_avalon",
		"crafting_total", "crafting_royal", "crafting_outlands", "crafting_avalon",
		"fishing_fame", "farming_fame", "crystal_league_fame",
	}

	assignments := make(map[string]interface{}, len(always)+len(statCols))
	for k, v := range always {
		assignments[k] = v
	}
	for _, col := range statCols {
		assignments[col] = gorm.Expr(
			fmt.Sprintf(
				"CASE WHEN excluded.last_seen > player_state.last_seen THEN excluded.%s ELSE player_state.%s END",
				col, col,
			),
		)
	}

	return p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
	}).Create(&states).Error
}

func (p *Poller) bulkUpsertSkips(ctx context.Context, states []models.PlayerState) error {
	return p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"next_poll_at",
			"priority",
		}),
	}).Create(&states).Error
}

func (p *Poller) bulkUpsertFailures(ctx context.Context, states []models.PlayerState) error {
	return p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"error_count",
			"next_poll_at",
		}),
	}).Create(&states).Error
}

func (p *Poller) fetchPlayer(ctx context.Context, playerID string) (*playerResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/players/%s", p.cfg.APIBase, playerID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("guid", randomGUID())
	q.Set("t", fmt.Sprintf("%d", time.Now().UnixNano()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if p.cfg.UserAgent != "" {
		req.Header.Set("User-Agent", p.cfg.UserAgent)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var pr playerResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func randomGUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func nullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func nullableInt64(val int64) *int64 {
	return &val
}

func nullableFloat64(val float64) *float64 {
	return &val
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
	base := 15 * time.Minute
	max := 24 * time.Hour
	// exponential backoff capped
	shift := errorCount
	if shift > 6 {
		shift = 6
	}
	backoff := base * (1 << shift)
	if backoff > max {
		backoff = max
	}
	return backoff
}

type playerResponse struct {
	Name               string        `json:"Name"`
	Id                 string        `json:"Id"`
	GuildName          string        `json:"GuildName"`
	GuildId            string        `json:"GuildId"`
	AllianceName       string        `json:"AllianceName"`
	AllianceId         string        `json:"AllianceId"`
	AllianceTag        string        `json:"AllianceTag"`
	KillFame           int64         `json:"KillFame"`
	DeathFame          int64         `json:"DeathFame"`
	FameRatio          float64       `json:"FameRatio"`
	Crafting           craftingStats `json:"Crafting,omitempty"`
	FishingFame        int64         `json:"FishingFame"`
	FarmingFame        int64         `json:"FarmingFame"`
	CrystalLeague      int64         `json:"CrystalLeague"`
	LifetimeStatistics lifetimeStats `json:"LifetimeStatistics"`
}

type craftingStats struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

type lifetimeStats struct {
	PvE       pveStats       `json:"PvE"`
	Gathering gatheringStats `json:"Gathering"`
	Crafting  craftingStats  `json:"Crafting"`
	Timestamp *time.Time     `json:"Timestamp"`
	Corrupted int64          `json:"CorruptedDungeon,omitempty"`
}

type pveStats struct {
	Total            int64 `json:"Total"`
	Royal            int64 `json:"Royal"`
	Outlands         int64 `json:"Outlands"`
	Avalon           int64 `json:"Avalon"`
	Hellgate         int64 `json:"Hellgate"`
	CorruptedDungeon int64 `json:"CorruptedDungeon"`
	Mists            int64 `json:"Mists"`
}

type gatheringStats struct {
	Fiber gatheringSplit `json:"Fiber"`
	Hide  gatheringSplit `json:"Hide"`
	Ore   gatheringSplit `json:"Ore"`
	Rock  gatheringSplit `json:"Rock"`
	Wood  gatheringSplit `json:"Wood"`
	All   gatheringSplit `json:"All"`
}

type gatheringSplit struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}
