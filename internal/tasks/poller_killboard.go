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

	"albionstats/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (p *KillboardPoller) collectPlayers(events []event, acc map[string]models.PlayerState) {
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

func (p *KillboardPoller) upsertPlayers(ctx context.Context, players map[string]models.PlayerState) error {
	batch := make([]models.PlayerState, 0, len(players))
	for _, pl := range players {
		player := pl // copy to avoid reference issues
		batch = append(batch, player)
	}

	condition := "COALESCE(player_state.last_polled, '-infinity') <= NOW() - interval '6 hours'"

	assignments := map[string]interface{}{
		"name":          gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.name ELSE player_state.name END", condition)),
		"guild_id":      gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.guild_id ELSE player_state.guild_id END", condition)),
		"guild_name":    gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.guild_name ELSE player_state.guild_name END", condition)),
		"alliance_id":   gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_id ELSE player_state.alliance_id END", condition)),
		"alliance_name": gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_name ELSE player_state.alliance_name END", condition)),
		"alliance_tag":  gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.alliance_tag ELSE player_state.alliance_tag END", condition)),
		"last_seen":     gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.last_seen ELSE player_state.last_seen END", condition)),
		"priority":      gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN 100 ELSE player_state.priority END", condition)),
		"next_poll_at":  gorm.Expr(fmt.Sprintf("CASE WHEN %s THEN excluded.next_poll_at ELSE player_state.next_poll_at END", condition)),
	}

	return p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "region"}, {Name: "player_id"}},
		DoUpdates: clause.Assignments(assignments),
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

func killboardNullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}
