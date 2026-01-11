package api

import (
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

type Client struct {
	httpClient *http.Client
	userAgent  string
}

func regionToBaseURL(region string) (string, error) {
	switch region {
	case "americas":
		return "https://gameinfo.albiononline.com", nil
	case "europe":
		return "https://gameinfo-ams.albiononline.com", nil
	case "asia":
		return "https://gameinfo-sgp.albiononline.com", nil
	default:
		return "", fmt.Errorf("invalid region: %s", region)
	}
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		userAgent: "AlbionStats-KillboardPoller/1.0",
	}
}

func (c *Client) FetchPlayer(region string, playerID string) (*PlayerResponse, error) {
	baseUrl, err := regionToBaseURL(region)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("%s/api/gameinfo/players/%s", baseUrl, playerID))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("guid", generateRandomGUID())
	q.Set("t", fmt.Sprintf("%d", time.Now().UnixNano()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
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

func (c *Client) FetchEvents(region string, limit int, offset int) ([]Event, error) {
	baseUrl, err := regionToBaseURL(region)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(baseUrl + "/api/gameinfo/events")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("guid", generateRandomGUID())
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
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

func generateRandomGUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

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

type CraftingStats struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

type LifetimeStats struct {
	PvE       PveStats       `json:"PvE"`
	Gathering GatheringStats `json:"Gathering"`
	Crafting  CraftingStats  `json:"Crafting"`
	Timestamp *time.Time     `json:"Timestamp"`
	Corrupted int64          `json:"CorruptedDungeon,omitempty"`
}

type PveStats struct {
	Total            int64 `json:"Total"`
	Royal            int64 `json:"Royal"`
	Outlands         int64 `json:"Outlands"`
	Avalon           int64 `json:"Avalon"`
	Hellgate         int64 `json:"Hellgate"`
	CorruptedDungeon int64 `json:"CorruptedDungeon"`
	Mists            int64 `json:"Mists"`
}

type GatheringStats struct {
	Fiber GatheringSplit `json:"Fiber"`
	Hide  GatheringSplit `json:"Hide"`
	Ore   GatheringSplit `json:"Ore"`
	Rock  GatheringSplit `json:"Rock"`
	Wood  GatheringSplit `json:"Wood"`
	All   GatheringSplit `json:"All"`
}

type GatheringSplit struct {
	Total    int64 `json:"Total"`
	Royal    int64 `json:"Royal"`
	Outlands int64 `json:"Outlands"`
	Avalon   int64 `json:"Avalon"`
}

type Event struct {
	EventID      int64         `json:"EventId"`
	TimeStamp    time.Time     `json:"TimeStamp"`
	Killer       Participant   `json:"Killer"`
	Victim       Participant   `json:"Victim"`
	Participants []Participant `json:"Participants"`
	GroupMembers []Participant `json:"GroupMembers"`
}

type Participant struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	GuildID      string `json:"GuildId"`
	GuildName    string `json:"GuildName"`
	AllianceID   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	AllianceTag  string `json:"AllianceTag"`
}
