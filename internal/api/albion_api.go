package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

// Client handles all Albion API interactions
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// NewClient creates a new Albion API client
func NewClient(baseURL, userAgent string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		userAgent:  userAgent,
	}
}

// FetchPlayer retrieves player data from the Albion API
func (c *Client) FetchPlayer(ctx context.Context, playerID string) (*PlayerResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/players/%s", c.baseURL, playerID))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("guid", generateRandomGUID())
	q.Set("t", fmt.Sprintf("%d", time.Now().UnixNano()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
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

	var pr PlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

// FetchEvents retrieves killboard events from the Albion API
func (c *Client) FetchEvents(ctx context.Context, limit, offset int) ([]Event, error) {
	u, err := url.Parse(c.baseURL + "/events")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("guid", generateRandomGUID())
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

// generateRandomGUID generates a random GUID for API requests
func generateRandomGUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// Fallback to timestamp-based value if entropy fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// NullableString converts an empty string to nil pointer
func NullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

// NullableInt64 converts an int64 to a pointer
func NullableInt64(val int64) *int64 {
	return &val
}

// NullableFloat64 converts a float64 to a pointer
func NullableFloat64(val float64) *float64 {
	return &val
}

// PlayerResponse represents the response from the Albion player API
type PlayerResponse struct {
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
	Crafting           CraftingStats `json:"Crafting,omitempty"`
	FishingFame        int64         `json:"FishingFame"`
	FarmingFame        int64         `json:"FarmingFame"`
	CrystalLeague      int64         `json:"CrystalLeague"`
	LifetimeStatistics LifetimeStats `json:"LifetimeStatistics"`
}

// CraftingStats represents crafting statistics
type CraftingStats struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

// LifetimeStats represents lifetime player statistics
type LifetimeStats struct {
	PvE       PveStats       `json:"PvE"`
	Gathering GatheringStats `json:"Gathering"`
	Crafting  CraftingStats  `json:"Crafting"`
	Timestamp *time.Time     `json:"Timestamp"`
	Corrupted int64          `json:"CorruptedDungeon,omitempty"`
}

// PveStats represents PvE statistics
type PveStats struct {
	Total            int64 `json:"Total"`
	Royal            int64 `json:"Royal"`
	Outlands         int64 `json:"Outlands"`
	Avalon           int64 `json:"Avalon"`
	Hellgate         int64 `json:"Hellgate"`
	CorruptedDungeon int64 `json:"CorruptedDungeon"`
	Mists            int64 `json:"Mists"`
}

// GatheringStats represents gathering statistics
type GatheringStats struct {
	Fiber GatheringSplit `json:"Fiber"`
	Hide  GatheringSplit `json:"Hide"`
	Ore   GatheringSplit `json:"Ore"`
	Rock  GatheringSplit `json:"Rock"`
	Wood  GatheringSplit `json:"Wood"`
	All   GatheringSplit `json:"All"`
}

// GatheringSplit represents gathering statistics split by region
type GatheringSplit struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

// Event represents a killboard event
type Event struct {
	EventID      int64         `json:"EventId"`
	TimeStamp    time.Time     `json:"TimeStamp"`
	Killer       Participant   `json:"Killer"`
	Victim       Participant   `json:"Victim"`
	Participants []Participant `json:"Participants"`
	GroupMembers []Participant `json:"GroupMembers"`
}

// Participant represents a participant in a killboard event
type Participant struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	GuildID      string `json:"GuildId"`
	GuildName    string `json:"GuildName"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	AllianceTag  string `json:"AllianceTag"`
}
