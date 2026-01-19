package tasks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

var (
	regionLimiters     map[string]*rate.Limiter
	regionLimitersOnce sync.Once
)

func getRegionLimiter(region string) *rate.Limiter {
	regionLimitersOnce.Do(func() {
		regionLimiters = map[string]*rate.Limiter{
			"americas": rate.NewLimiter(rate.Limit(4), 4),
			"europe":   rate.NewLimiter(rate.Limit(4), 4),
			"asia":     rate.NewLimiter(rate.Limit(4), 4),
		}
	})
	return regionLimiters[region]
}

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

	if limiter := getRegionLimiter(region); limiter != nil {
		if err := limiter.Wait(context.Background()); err != nil {
			return nil, err
		}
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

	if limiter := getRegionLimiter(region); limiter != nil {
		if err := limiter.Wait(context.Background()); err != nil {
			return nil, err
		}
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

func (c *Client) FetchBattles(region string, offset, limit int) (BattlesResponse, error) {
	baseUrl, err := regionToBaseURL(region)
	if err != nil {
		return nil, err
	}

	if limiter := getRegionLimiter(region); limiter != nil {
		if err := limiter.Wait(context.Background()); err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(fmt.Sprintf("%s/api/gameinfo/battles", baseUrl))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("sort", "recent")
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var battles BattlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&battles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return battles, nil
}

func (c *Client) FetchBattleEvents(region string, battleID int64, offset, limit int) ([]Event, error) {
	baseUrl, err := regionToBaseURL(region)
	if err != nil {
		return nil, err
	}

	if limiter := getRegionLimiter(region); limiter != nil {
		if err := limiter.Wait(context.Background()); err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(fmt.Sprintf("%s/api/gameinfo/events/battle/%d", baseUrl, battleID))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("limit", fmt.Sprintf("%d", limit))
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
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
	Total    int64 `json:"Total,omitempty"`
	Royal    int64 `json:"Royal,omitempty"`
	Outlands int64 `json:"Outlands,omitempty"`
	Avalon   int64 `json:"Avalon,omitempty"`
}

type LifetimeStats struct {
	PvE          PveStats       `json:"PvE,omitempty"`
	Gathering    GatheringStats `json:"Gathering,omitempty"`
	Crafting     CraftingStats  `json:"Crafting,omitempty"`
	CrystalLeague int64         `json:"CrystalLeague,omitempty"`
	FishingFame   int64         `json:"FishingFame,omitempty"`
	FarmingFame   int64         `json:"FarmingFame,omitempty"`
	Timestamp     *time.Time    `json:"Timestamp,omitempty"`
}

type PveStats struct {
	Total            int64 `json:"Total,omitempty"`
	Royal            int64 `json:"Royal,omitempty"`
	Outlands         int64 `json:"Outlands,omitempty"`
	Avalon           int64 `json:"Avalon,omitempty"`
	Hellgate         int64 `json:"Hellgate,omitempty"`
	CorruptedDungeon int64 `json:"CorruptedDungeon,omitempty"`
	Mists            int64 `json:"Mists,omitempty"`
}

type GatheringStats struct {
	Fiber GatheringSplit `json:"Fiber,omitempty"`
	Hide  GatheringSplit `json:"Hide,omitempty"`
	Ore   GatheringSplit `json:"Ore,omitempty"`
	Rock  GatheringSplit `json:"Rock,omitempty"`
	Wood  GatheringSplit `json:"Wood,omitempty"`
	All   GatheringSplit `json:"All,omitempty"`
}

type GatheringSplit struct {
	Total    int64 `json:"Total,omitempty"`
	Royal    int64 `json:"Royal,omitempty"`
	Outlands int64 `json:"Outlands,omitempty"`
	Avalon   int64 `json:"Avalon,omitempty"`
}

type EquipmentItem struct {
	Type          string         `json:"Type"`
	Count         int32          `json:"Count"`
	Quality       int32          `json:"Quality"`
	ActiveSpells  []interface{}  `json:"ActiveSpells,omitempty"`
	PassiveSpells []interface{}  `json:"PassiveSpells,omitempty"`
	LegendarySoul *LegendarySoul `json:"LegendarySoul,omitempty"`
}

type LegendarySoul struct {
	ID                      string    `json:"id"`
	Subtype                 int32     `json:"subtype"`
	Era                     int32     `json:"era"`
	Name                    *string   `json:"name,omitempty"`
	LastEquipped            time.Time `json:"lastEquipped"`
	AttunedPlayer           string    `json:"attunedPlayer"`
	AttunedPlayerName       string    `json:"attunedPlayerName"`
	Attunement              int64     `json:"attunement"`
	AttunementSpentSinceReset int64   `json:"attunementSpentSinceReset"`
	AttunementSpent         int64     `json:"attunementSpent"`
	Quality                 int32     `json:"quality"`
	CraftedBy               string    `json:"craftedBy"`
	Traits                  []Trait   `json:"traits"`
	PvPFameGained           int64     `json:"PvPFameGained"`
}

type Trait struct {
	Roll         float64       `json:"roll"`
	PendingRolls []interface{} `json:"pendingRolls,omitempty"`
	PendingTraits []interface{} `json:"pendingTraits,omitempty"`
	Value        float64       `json:"value"`
	Trait        string        `json:"trait"`
	MinValue     float64       `json:"minvalue"`
	MaxValue     float64       `json:"maxvalue"`
}

type Event struct {
	GroupMemberCount     int32         `json:"groupMemberCount,omitempty"`
	NumberOfParticipants int32         `json:"numberOfParticipants,omitempty"`
	EventID              int64         `json:"EventId"`
	TimeStamp            time.Time     `json:"TimeStamp"`
	Version              int32         `json:"Version,omitempty"`
	Killer               Participant   `json:"Killer"`
	Victim               Participant   `json:"Victim"`
	TotalVictimKillFame  int64         `json:"TotalVictimKillFame,omitempty"`
	Location             interface{}   `json:"Location,omitempty"`
	Participants         []Participant `json:"Participants"`
	GroupMembers         []Participant `json:"GroupMembers"`
	GvGMatch             interface{}   `json:"GvGMatch,omitempty"`
	BattleID             int64         `json:"BattleId,omitempty"`
	KillArea             string        `json:"KillArea,omitempty"`
	Category             interface{}   `json:"Category,omitempty"`
	Type                 string        `json:"Type,omitempty"`
}

type Participant struct {
	AverageItemPower     float64                `json:"AverageItemPower,omitempty"`
	Equipment            map[string]*EquipmentItem `json:"Equipment,omitempty"`
	Inventory            []interface{}          `json:"Inventory,omitempty"`
	Name                 string                 `json:"Name"`
	ID                   string                 `json:"Id"`
	GuildName            string                 `json:"GuildName"`
	GuildID              string                 `json:"GuildId"`
	AllianceName         string                 `json:"AllianceName"`
	AllianceID           string                 `json:"AllianceId"`
	AllianceTag          string                 `json:"AllianceTag"`
	Avatar               string                 `json:"Avatar,omitempty"`
	AvatarRing           string                 `json:"AvatarRing,omitempty"`
	DeathFame            int64                  `json:"DeathFame,omitempty"`
	KillFame             int64                  `json:"KillFame,omitempty"`
	FameRatio            float64                `json:"FameRatio,omitempty"`
	LifetimeStatistics   *LifetimeStats         `json:"LifetimeStatistics,omitempty"`
	DamageDone           float64                `json:"DamageDone,omitempty"`
	SupportHealingDone   float64                `json:"SupportHealingDone,omitempty"`
}

type BattlePlayer struct {
	Name         string `json:"name"`
	Kills        int32  `json:"kills"`
	Deaths       int32  `json:"deaths"`
	KillFame     int64  `json:"killFame"`
	GuildName    *string `json:"guildName"`
	GuildID      *string `json:"guildId"`
	AllianceName *string `json:"allianceName"`
	AllianceID   *string `json:"allianceId"`
	ID           string `json:"id"`
}

type BattleGuild struct {
	Name       string `json:"name"`
	Kills      int32  `json:"kills"`
	Deaths     int32  `json:"deaths"`
	KillFame   int64  `json:"killFame"`
	Alliance   *string `json:"alliance"`
	AllianceID *string `json:"allianceId"`
	ID         *string `json:"id"`
}

type BattleAlliance struct {
	Name     string `json:"name"`
	Kills    int32  `json:"kills"`
	Deaths   int32  `json:"deaths"`
	KillFame int64  `json:"killFame"`
	ID       string `json:"id"`
}

type Battle struct {
	ID          int64                      `json:"id"`
	StartTime   time.Time                  `json:"startTime"`
	EndTime     time.Time                  `json:"endTime,omitempty"`
	Timeout     time.Time                  `json:"timeout,omitempty"`
	TotalFame   int64                      `json:"totalFame"`
	TotalKills  int32                      `json:"totalKills"`
	ClusterName *string                    `json:"clusterName,omitempty"`
	Players     map[string]BattlePlayer    `json:"players"`
	Guilds      map[string]BattleGuild     `json:"guilds"`
	Alliances   map[string]BattleAlliance  `json:"alliances"`
	BattleTimeout int32                    `json:"battle_TIMEOUT"`
}

type BattlesResponse []Battle