package killboard_poller

import (
	"log/slog"
	"time"

	"albionstats/internal/api"
	"albionstats/internal/sqlite"
)

type Config struct {
	APIClient      *api.Client
	SQLite         *sqlite.SQLite
	Logger         *slog.Logger
	Region         string
	PageSize       int
	MaxPages       int
	EventsInterval time.Duration
}

type KillboardPoller struct {
	apiClient      *api.Client
	sqlite         *sqlite.SQLite
	log            *slog.Logger
	eventsInterval time.Duration
	pageSize       int
	region         string
}

func NewKillboardPoller(cfg Config) (*KillboardPoller, error) {
	return &KillboardPoller{
		apiClient:      cfg.APIClient,
		sqlite:         cfg.SQLite,
		log:            cfg.Logger.With("component", "killboard_poller", "region", cfg.Region),
		eventsInterval: cfg.EventsInterval,
		pageSize:       cfg.PageSize,
		region:         cfg.Region,
	}, nil
}

func (p *KillboardPoller) Run() {
	p.log.Info("killboard polling started", "interval", p.eventsInterval, "page_size", p.pageSize)

	ticker := time.NewTicker(p.eventsInterval)
	defer ticker.Stop()

	p.runBatch() // Run once immediately

	for range ticker.C {
		p.runBatch()
	}
}

func (p *KillboardPoller) runBatch() {
	p.log.Info("fetch killboard events", "limit", p.pageSize, "offset", 0)
	events, err := p.apiClient.FetchEvents(p.region, p.pageSize, 0)
	if err != nil {
		p.log.Warn("fetch killboard events failed", "err", err)
		return
	}

	if len(events) == 0 {
		p.log.Warn("no events returned")
		return
	}

	playerMap := make(map[string]sqlite.PlayerPoll)
	p.collectPlayers(events, playerMap)

	if len(playerMap) == 0 {
		return
	}

	if err := p.sqlite.UpsertPlayerPolls(playerMap); err != nil {
		p.log.Error("upsert player polls failed", "err", err, "players", len(playerMap))
		return
	}
	p.log.Info("upserted player polls", "count", len(playerMap))
}

func (p *KillboardPoller) collectPlayers(events []api.Event, acc map[string]sqlite.PlayerPoll) {
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
			playerPoll := sqlite.PlayerPoll{
				Region:                p.region,
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
