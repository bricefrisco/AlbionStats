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
	DBDSN               string
	APIBase             string
	PageSize            int
	MaxPages            int
	EventsInterval      time.Duration
	KillboardPageSize   int
	KillboardMaxPages   int
	KillboardInterval   time.Duration
	HTTPTimeout         time.Duration
	PlayerRate          int
	PlayerBatch         int
	UserAgent           string
	APIPort             string
}

const (
	defaultPageSize            = 50
	defaultMaxPages            = 1
	defaultEventsInterval      = 10 * time.Second
	defaultKillboardPageSize   = 51
	defaultKillboardMaxPages   = 1
	defaultKillboardInterval   = 60 * time.Second
	defaultPlayerRate          = 6
	defaultPlayerBatch         = 100
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
		DBDSN:               valueWithDefault(values, "ALBION_DB_DSN", defaultDBDSN),
		PageSize:            intFrom(values, "ALBION_EVENTS_PAGE_SIZE", defaultPageSize),
		MaxPages:            intFrom(values, "ALBION_EVENTS_MAX_PAGES", defaultMaxPages),
		EventsInterval:      durationFrom(values, "ALBION_EVENTS_INTERVAL", defaultEventsInterval),
		KillboardPageSize:   intFrom(values, "ALBION_KILLBOARD_PAGE_SIZE", defaultKillboardPageSize),
		KillboardMaxPages:   intFrom(values, "ALBION_KILLBOARD_MAX_PAGES", defaultKillboardMaxPages),
		KillboardInterval:   durationFrom(values, "ALBION_KILLBOARD_INTERVAL", defaultKillboardInterval),
		PlayerRate:          intFrom(values, "ALBION_PLAYER_RATE", defaultPlayerRate),
		PlayerBatch:         intFrom(values, "ALBION_PLAYER_BATCH", defaultPlayerBatch),
		APIPort:             valueWithDefault(values, "API_PORT", defaultAPIPort),
	}

	if cfg.PageSize <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_PAGE_SIZE: %d", cfg.PageSize)
	}
	if cfg.MaxPages <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_MAX_PAGES: %d", cfg.MaxPages)
	}
	if cfg.EventsInterval <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_EVENTS_INTERVAL: %v", cfg.EventsInterval)
	}
	if cfg.PlayerRate <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_PLAYER_RATE: %d", cfg.PlayerRate)
	}
	if cfg.PlayerBatch <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_PLAYER_BATCH: %d", cfg.PlayerBatch)
	}
	if cfg.KillboardPageSize <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_KILLBOARD_PAGE_SIZE: %d", cfg.KillboardPageSize)
	}
	if cfg.KillboardMaxPages <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_KILLBOARD_MAX_PAGES: %d", cfg.KillboardMaxPages)
	}
	if cfg.KillboardInterval <= 0 {
		return Config{}, fmt.Errorf("invalid ALBION_KILLBOARD_INTERVAL: %v", cfg.KillboardInterval)
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
