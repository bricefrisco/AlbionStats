package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration loaded from a config file.
type Config struct {
	DBDSN          string
	APIBase        string
	Region         string
	PageSize       int
	MaxPages       int
	EventsInterval time.Duration
	HTTPTimeout    time.Duration
	PlayerRate     int
	PlayerBatch    int
	UserAgent      string
}

const (
	defaultAPIBase        = "https://gameinfo.albiononline.com/api/gameinfo"
	defaultRegion         = "americas"
	defaultPageSize       = 50
	defaultMaxPages       = 1
	defaultEventsInterval = 10 * time.Second
	defaultHTTPTimeout    = 10 * time.Second
	defaultPlayerRate     = 6
	defaultPlayerBatch    = 100
	defaultUserAgent      = "AlbionStats-KillboardPoller/1.0"
	defaultConfigPath     = ".env"
)

// Load reads configuration from a simple KEY=VALUE file (default .env).
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
		DBDSN:          strings.TrimSpace(values["ALBION_DB_DSN"]),
		APIBase:        valueWithDefault(values, "ALBION_API_BASE", defaultAPIBase),
		Region:         strings.ToLower(valueWithDefault(values, "ALBION_REGION", defaultRegion)),
		PageSize:       intFrom(values, "ALBION_EVENTS_PAGE_SIZE", defaultPageSize),
		MaxPages:       intFrom(values, "ALBION_EVENTS_MAX_PAGES", defaultMaxPages),
		EventsInterval: durationFrom(values, "ALBION_EVENTS_INTERVAL", defaultEventsInterval),
		HTTPTimeout:    durationFrom(values, "ALBION_HTTP_TIMEOUT", defaultHTTPTimeout),
		PlayerRate:     intFrom(values, "ALBION_PLAYER_RATE", defaultPlayerRate),
		PlayerBatch:    intFrom(values, "ALBION_PLAYER_BATCH", defaultPlayerBatch),
		UserAgent:      valueWithDefault(values, "ALBION_USER_AGENT", defaultUserAgent),
	}

	if cfg.DBDSN == "" {
		return Config{}, errors.New("ALBION_DB_DSN is required")
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
