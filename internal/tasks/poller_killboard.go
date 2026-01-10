package tasks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"albionstats/internal/api"
	"albionstats/internal/database"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type KillboardConfig struct {
	PageSize       int
	MaxPages       int
	EventsInterval time.Duration
	Region         string
	UserAgent      string
	HTTPTimeout    time.Duration
}

type KillboardPoller struct {
	apiClient *api.Client
	db        *gorm.DB
	cfg       KillboardConfig
	log       *slog.Logger
}

func NewKillboardPoller(db *gorm.DB, logger *slog.Logger, cfg KillboardConfig) (*KillboardPoller, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	apiClient := api.NewClient(cfg.Region, cfg.UserAgent, cfg.HTTPTimeout)
	return &KillboardPoller{
		apiClient: apiClient,
		db:        db,
		cfg:       cfg,
		log:       logger.With("component", "killboard_poller", "region", cfg.Region),
	}, nil
}

// Run fetches events pages and upserts discovered players into the database periodically.
func (p *KillboardPoller) Run(ctx context.Context) {
	p.log.Info("killboard polling started", "interval", p.cfg.EventsInterval, "page_size", p.cfg.PageSize)

	ticker := time.NewTicker(p.cfg.EventsInterval)
	defer ticker.Stop()

	// Run once immediately
	p.runBatch(ctx)

	for {
		select {
		case <-ctx.Done():
			p.log.Info("killboard polling stopped")
			return
		case <-ticker.C:
			p.runBatch(ctx)
		}
	}
}

func (p *KillboardPoller) runBatch(ctx context.Context) {
	start := time.Now()
	p.log.Info("fetch events", "limit", p.cfg.PageSize, "offset", 0)
	playerMap := make(map[string]database.PlayerPoll)
	events, err := p.apiClient.FetchEvents(ctx, p.cfg.PageSize, 0)
	if err != nil {
		p.log.Warn("fetch events failed", "err", err)
		return
	}
	if len(events) == 0 {
		p.log.Warn("no events returned")
		return
	}

	p.collectPlayers(events, playerMap)
	p.log.Info("events processed", "events", len(events), "players_collected", len(playerMap))

	if len(playerMap) == 0 {
		return
	}

	if err := database.UpsertKillboardPlayerPolls(ctx, p.db, playerMap); err != nil {
		p.log.Error("upsert player polls failed", "err", err, "players", len(playerMap))
		return
	}
	p.log.Info("upserted player polls", "count", len(playerMap), "duration_ms", time.Since(start).Milliseconds())
}

func (p *KillboardPoller) collectPlayers(events []api.Event, acc map[string]database.PlayerPoll) {
	now := time.Now().UTC()
	for _, ev := range events {
		add := func(participant api.Participant) {
			if participant.ID == "" {
				return
			}
			key := participant.ID
			if _, exists := acc[key]; exists {
				return
			}
			playerPoll := database.PlayerPoll{
				Region:                p.cfg.Region,
				PlayerID:              participant.ID,
				NextPollAt:            now,
				KillboardLastActivity: &ev.TimeStamp,
			}
			acc[key] = playerPoll
		}

		add(ev.Killer)
		add(ev.Victim)
		for _, part := range ev.Participants {
			add(part)
		}
		for _, gm := range ev.GroupMembers {
			add(gm)
		}
	}
}
