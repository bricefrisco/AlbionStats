package tasks

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type CollectorConfig struct {
	Interval time.Duration
}

type Collector struct {
	db     *gorm.DB
	config CollectorConfig
}

func NewCollector(db *gorm.DB, config CollectorConfig) *Collector {
	return &Collector{
		db:     db,
		config: config,
	}
}

func (c *Collector) Start(ctx context.Context) {
	log.Printf("metrics: starting collector with interval %v", c.config.Interval)

	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	// Run once immediately
	c.collect(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("metrics: collector stopped")
			return
		case <-ticker.C:
			c.collect(ctx)
		}
	}
}

func (c *Collector) collect(ctx context.Context) {
	log.Printf("metrics: collecting metrics")

	// Execute the metrics insertion query
	err := c.db.WithContext(ctx).Exec(`
		INSERT INTO metrics_timeseries (metric, ts, value)
		VALUES
			('players_total', now(), (SELECT COUNT(*) FROM player_state)),
			('snapshots_estimated', now(),
				(SELECT approximate_row_count('player_stats_snapshots')))
	`).Error

	if err != nil {
		log.Printf("metrics: failed to collect metrics: %v", err)
		return
	}

	log.Printf("metrics: successfully collected metrics")
}
