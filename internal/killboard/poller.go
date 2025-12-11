package killboard

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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	APIBase        string
	PageSize       int
	MaxPages       int
	EventsInterval time.Duration
	Region         string
	UserAgent      string
}

type Poller struct {
	client *http.Client
	db     *gorm.DB
	cfg    Config
}

func New(client *http.Client, db *gorm.DB, cfg Config) *Poller {
	return &Poller{
		client: client,
		db:     db,
		cfg:    cfg,
	}
}

// Run fetches events pages and upserts discovered players into the database.
func (p *Poller) Run(ctx context.Context) error {
	log.Printf("poller: fetch events limit=%d offset=0", p.cfg.PageSize)
	playerMap := make(map[string]models.PlayerState)
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

	if err := p.upsertPlayers(ctx, playerMap); err != nil {
		return err
	}
	log.Printf("poller: upserted players=%d", len(playerMap))
	return nil
}

func (p *Poller) fetchEvents(ctx context.Context, limit, offset int) ([]event, error) {
	u, err := url.Parse(p.cfg.APIBase + "/events")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("guid", randomGUID())
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

func randomGUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// Fallback to timestamp-based value if entropy fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func (p *Poller) collectPlayers(events []event, acc map[string]models.PlayerState) {
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
			player := models.PlayerState{
				Region:       p.cfg.Region,
				PlayerID:     participant.ID,
				Name:         participant.Name,
				Priority:     100,
				NextPollAt:   now,
				ErrorCount:   0,
				LastSeen:     &lastSeen,
				GuildID:      nullableString(participant.GuildID),
				GuildName:    nullableString(participant.GuildName),
				AllianceID:   nullableString(participant.AllianceID),
				AllianceName: nullableString(participant.AllianceName),
				AllianceTag:  nullableString(participant.AllianceTag),
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

func (p *Poller) upsertPlayers(ctx context.Context, players map[string]models.PlayerState) error {
	batch := make([]models.PlayerState, 0, len(players))
	for _, pl := range players {
		player := pl // copy to avoid reference issues
		batch = append(batch, player)
	}

	return p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":          clause.Column{Table: "excluded", Name: "name"},
			"guild_id":      clause.Column{Table: "excluded", Name: "guild_id"},
			"guild_name":    clause.Column{Table: "excluded", Name: "guild_name"},
			"alliance_id":   clause.Column{Table: "excluded", Name: "alliance_id"},
			"alliance_name": clause.Column{Table: "excluded", Name: "alliance_name"},
			"alliance_tag":  clause.Column{Table: "excluded", Name: "alliance_tag"},
			"last_seen":     clause.Column{Table: "excluded", Name: "last_seen"},
			"priority":      100,
			"next_poll_at":  gorm.Expr("NOW()"),
		}),
	}).Create(&batch).Error
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

func nullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}
