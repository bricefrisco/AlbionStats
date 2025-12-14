package tasks

import (
	"context"
	"log"
	"time"

	"albionstats/internal/api"
	"albionstats/internal/database"

	"gorm.io/gorm"
)

type KillboardConfig struct {
	APIBase        string
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
}

func NewKillboardPoller(db *gorm.DB, cfg KillboardConfig) *KillboardPoller {
	apiClient := api.NewClient(cfg.APIBase, cfg.UserAgent, cfg.HTTPTimeout)
	return &KillboardPoller{
		apiClient: apiClient,
		db:        db,
		cfg:       cfg,
	}
}

// Run fetches events pages and upserts discovered players into the database periodically.
func (p *KillboardPoller) Run(ctx context.Context) {
	log.Printf("starting periodic killboard polling with interval %v", p.cfg.EventsInterval)

	ticker := time.NewTicker(p.cfg.EventsInterval)
	defer ticker.Stop()

	// Run once immediately
	p.runBatch(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopped")
			return
		case <-ticker.C:
			p.runBatch(ctx)
		}
	}
}

func (p *KillboardPoller) runBatch(ctx context.Context) {
	log.Printf("fetch events limit=%d offset=0", p.cfg.PageSize)
	playerMap := make(map[string]database.PlayerPoll)
	events, err := p.apiClient.FetchEvents(ctx, p.cfg.PageSize, 0)
	if err != nil {
		log.Printf("fetch events error: %v", err)
		return
	}
	if len(events) == 0 {
		log.Printf("no events returned")
		return
	}

	p.collectPlayers(events, playerMap)
	log.Printf("events=%d players_collected=%d", len(events), len(playerMap))

	if len(playerMap) == 0 {
		return
	}

	if err := database.UpsertPlayerPolls(ctx, p.db, playerMap); err != nil {
		log.Printf("upsert player polls error: %v", err)
		return
	}
	log.Printf("upserted player polls=%d", len(playerMap))
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
				ErrorCount:            0,
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
