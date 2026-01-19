package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBDSN                    string
	APIBase                  string
	EventsPageSize           int
	EventsMaxPages           int
	EventsInterval           time.Duration
	BattleboardPageSize      int
	BattleboardMaxPages      int
	BattleboardInterval      time.Duration
	HTTPTimeout              time.Duration
	PlayerBatch              int
	PlayerWorkerCount        int
	UserAgent                string
	APIPort                  string
}

const (
	defaultEventsPageSize      = 50
	defaultEventsMaxPages      = 1
	defaultEventsInterval      = 10 * time.Second
	defaultBattleboardPageSize = 51
	defaultBattleboardMaxPages = 1
	defaultBattleboardInterval = 60 * time.Second
	defaultPlayerBatch         = 100
	defaultPlayerWorkerCount   = 5
	defaultAPIPort             = "8080"
	defaultDBDSN               = "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	defaultConfigPath          = ".env"
)

func Load() (Config, error) {
	path := strings.TrimSpace(os.Getenv("ALBION_CONFIG_FILE"))
	if path == "" {
		path = defaultConfigPath
	}

	values, err := parseEnvFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("load config file %s: %w", path, err)
	}

	cfg := Config{
		DBDSN:                    valueWithDefault(values, "ALBION_DB_DSN", defaultDBDSN),
		EventsPageSize:           intFrom(values, "ALBION_EVENTS_PAGE_SIZE", defaultEventsPageSize),
		EventsMaxPages:           intFrom(values, "ALBION_EVENTS_MAX_PAGES", defaultEventsMaxPages),
		EventsInterval:           durationFrom(values, "ALBION_EVENTS_INTERVAL", defaultEventsInterval),
		BattleboardPageSize:      intFrom(values, "ALBION_BATTLE_BOARD_PAGE_SIZE", defaultBattleboardPageSize),
		BattleboardMaxPages:      intFrom(values, "ALBION_BATTLE_BOARD_MAX_PAGES", defaultBattleboardMaxPages),
		BattleboardInterval:      durationFrom(values, "ALBION_BATTLE_BOARD_INTERVAL", defaultBattleboardInterval),
		PlayerBatch:              intFrom(values, "ALBION_PLAYER_BATCH", defaultPlayerBatch),
		PlayerWorkerCount:        intFrom(values, "ALBION_PLAYER_WORKER_COUNT", defaultPlayerWorkerCount),
		APIPort:                  valueWithDefault(values, "API_PORT", defaultAPIPort),
	}

	if cfg.EventsPageSize <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_PAGE_SIZE: %d", cfg.EventsPageSize)
	}
	if cfg.EventsMaxPages <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_MAX_PAGES: %d", cfg.EventsMaxPages)
	}
	if cfg.EventsInterval <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_INTERVAL: %v", cfg.EventsInterval)
	}
	if cfg.BattleboardPageSize <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_BATTLE_BOARD_PAGE_SIZE: %d", cfg.BattleboardPageSize)
	}
	if cfg.BattleboardMaxPages <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_BATTLE_BOARD_MAX_PAGES: %d", cfg.BattleboardMaxPages)
	}
	if cfg.BattleboardInterval <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_BATTLE_BOARD_INTERVAL: %v", cfg.BattleboardInterval)
	}
	if cfg.PlayerBatch <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_PLAYER_BATCH: %d", cfg.PlayerBatch)
	}
	if cfg.PlayerWorkerCount <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_PLAYER_WORKER_COUNT: %d", cfg.PlayerWorkerCount)
	}

	cfg.APIBase = strings.TrimRight(cfg.APIBase, "/")
	return cfg, nil
}

func valueWithDefault(values map[string]string, key, def string) string {
	if val := strings.TrimSpace(values[key]); val != "" {
		return val
	}
	return def
}

func intFrom(values map[string]string, key string, def int) int {
	val := strings.TrimSpace(values[key])
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}

func durationFrom(values map[string]string, key string, def time.Duration) time.Duration {
	val := strings.TrimSpace(values[key])
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		values[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
