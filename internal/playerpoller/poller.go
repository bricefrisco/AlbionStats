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
	"time"

	"albionstats/internal/models"

	"errors"

	"gorm.io/gorm"
)

type Config struct {
	APIBase    string
	PageSize   int
	RatePerSec int
	UserAgent  string
}

type Poller struct {
	client *http.Client
	db     *gorm.DB
	cfg    Config
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

	log.Printf("player-poller: batch size=%d rate=%d/s", len(players), p.cfg.RatePerSec)

	for i, pl := range players {
		if pl.LastPolled != nil && time.Since(*pl.LastPolled) < 6*time.Hour {
			nextPoll := pl.LastPolled.Add(6 * time.Hour)
			if err := p.db.WithContext(ctx).Model(&models.PlayerState{}).
				Where("region = ? AND player_id = ?", pl.Region, pl.PlayerID).
				Updates(map[string]interface{}{
					"next_poll_at": nextPoll,
					"priority":     200,
				}).Error; err != nil {
				log.Printf("player-poller: player=%s skip-update err=%v", pl.PlayerID, err)
			} else {
				log.Printf("player-poller: player=%s last polled too recent; next_poll_at=%s", pl.PlayerID, nextPoll.Format(time.RFC3339))
			}
			continue
		}
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
		if err := p.processPlayer(ctx, pl); err != nil {
			log.Printf("player-poller: player=%s err=%v", pl.PlayerID, err)
			if err2 := p.handleFailure(ctx, pl); err2 != nil {
				log.Printf("player-poller: failure handling error for %s: %v", pl.PlayerID, err2)
			}
		} else {
			log.Printf("player-poller: player=%s ok", pl.PlayerID)
		}
	}

	return nil
}

func (p *Poller) processPlayer(ctx context.Context, pl models.PlayerState) error {
	ts := time.Now().UTC()
	resp, err := p.fetchPlayer(ctx, pl.PlayerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if delErr := p.db.WithContext(ctx).Delete(&models.PlayerState{}, "region = ? AND player_id = ?", pl.Region, pl.PlayerID).Error; delErr != nil {
				return fmt.Errorf("player %s 404 delete: %w", pl.PlayerID, delErr)
			}
			log.Printf("player-poller: removed player=%s due to 404", pl.PlayerID)
			return nil
		}
		return fmt.Errorf("fetch player %s: %w", pl.PlayerID, err)
	}

	apiTS := resp.LifetimeStatistics.Timestamp
	if apiTS == nil {
		// Remove player if API payload lacks timestamp (assumed invalid)
		if err := p.db.WithContext(ctx).Delete(&models.PlayerState{}, "region = ? AND player_id = ?", pl.Region, pl.PlayerID).Error; err != nil {
			return fmt.Errorf("player %s: delete on missing timestamp: %w", pl.PlayerID, err)
		}
		log.Printf("player-poller: removed player=%s due to missing API timestamp", pl.PlayerID)
		return nil
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

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("region = ? AND player_id = ?", pl.Region, pl.PlayerID).
			Save(&update).Error; err != nil {
			return err
		}
		if err := tx.Create(&snapshot).Error; err != nil {
			return err
		}
		return nil
	})
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

func (p *Poller) handleFailure(ctx context.Context, pl models.PlayerState) error {
	nextErrCount := pl.ErrorCount + 1
	backoff := failureBackoff(nextErrCount)
	nextPoll := time.Now().UTC().Add(backoff)
	return p.db.WithContext(ctx).Model(&models.PlayerState{}).
		Where("region = ? AND player_id = ?", pl.Region, pl.PlayerID).
		Updates(map[string]interface{}{
			"error_count":  nextErrCount,
			"next_poll_at": nextPoll,
		}).Error
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
