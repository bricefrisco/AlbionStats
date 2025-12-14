package tasks

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
	client *http.Client
	db     *gorm.DB
	cfg    KillboardConfig
}

func NewKillboardPoller(db *gorm.DB, cfg KillboardConfig) *KillboardPoller {
	client := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}
	return &KillboardPoller{
		client: client,
		db:     db,
		cfg:    cfg,
	}
}

// Run fetches events pages and upserts discovered players into the database.
func (p *KillboardPoller) Run(ctx context.Context) error {
	log.Printf("poller: fetch events limit=%d offset=0", p.cfg.PageSize)
	playerMap := make(map[string]database.PlayerState)
	events, err := p.fetchEvents(ctx, p.cfg.PageSize, 0)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}
	if len(events) == 0 {
		log.Printf("poller: no events returned")
		return nil
	}

	p.collectPlayers(events, playerMap)
	log.Printf("poller: events=%d players_collected=%d", len(events), len(playerMap))

	if len(playerMap) == 0 {
		return nil
	}

	if err := database.UpsertPlayers(ctx, p.db, playerMap); err != nil {
		return err
	}
	log.Printf("poller: upserted players=%d", len(playerMap))
	return nil
}

func (p *KillboardPoller) fetchEvents(ctx context.Context, limit, offset int) ([]event, error) {
	u, err := url.Parse(p.cfg.APIBase + "/events")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("guid", killboardRandomGUID())
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var events []event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

func killboardRandomGUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// Fallback to timestamp-based value if entropy fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func (p *KillboardPoller) collectPlayers(events []event, acc map[string]database.PlayerState) {
	now := time.Now().UTC()
	for _, ev := range events {
		lastSeen := ev.TimeStamp
		add := func(participant participant) {
			if participant.ID == "" {
				return
			}
			key := participant.ID
			if _, exists := acc[key]; exists {
				return
			}
			player := database.PlayerState{
				Region:       p.cfg.Region,
				PlayerID:     participant.ID,
				Name:         participant.Name,
				Priority:     100,
				NextPollAt:   now,
				ErrorCount:   0,
				LastSeen:     &lastSeen,
				GuildID:      killboardNullableString(participant.GuildID),
				GuildName:    killboardNullableString(participant.GuildName),
				AllianceID:   killboardNullableString(participant.AllianceID),
				AllianceName: killboardNullableString(participant.AllianceName),
				AllianceTag:  killboardNullableString(participant.AllianceTag),
			}
			acc[key] = player
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

type event struct {
	EventID      int64         `json:"EventId"`
	TimeStamp    time.Time     `json:"TimeStamp"`
	Killer       participant   `json:"Killer"`
	Victim       participant   `json:"Victim"`
	Participants []participant `json:"Participants"`
	GroupMembers []participant `json:"GroupMembers"`
}

type participant struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	GuildID      string `json:"GuildId"`
	GuildName    string `json:"GuildName"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	AllianceTag  string `json:"AllianceTag"`
}

func killboardNullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}
