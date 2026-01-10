package tasks

// import (
// 	"context"
// 	"fmt"
// 	"log/slog"
// 	"time"

// 	"gorm.io/gorm"
// )

// type CollectorConfig struct {
// 	Interval time.Duration
// }

// type Collector struct {
// 	db     *gorm.DB
// 	config CollectorConfig
// 	log    *slog.Logger
// }

// func NewCollector(db *gorm.DB, logger *slog.Logger, config CollectorConfig) (*Collector, error) {
// 	if logger == nil {
// 		return nil, fmt.Errorf("logger is required")
// 	}

// 	return &Collector{
// 		db:     db,
// 		config: config,
// 		log:    logger.With("component", "metrics_collector"),
// 	}, nil
// }

// func (c *Collector) Run(ctx context.Context) {
// 	c.log.Info("collector started", "interval", c.config.Interval)

// 	ticker := time.NewTicker(c.config.Interval)
// 	defer ticker.Stop()

// 	// Run once immediately
// 	c.collect(ctx)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			c.log.Info("collector stopped")
// 			return
// 		case <-ticker.C:
// 			c.collect(ctx)
// 		}
// 	}
// }

// func (c *Collector) collect(ctx context.Context) {
// 	start := time.Now()
// 	c.log.Info("metrics collection started")

// 	// Execute the metrics insertion query
// 	err := c.db.WithContext(ctx).Exec(`
// 		INSERT INTO metrics (metric, ts, value)
// 		VALUES
// 			('players_total', now(), (SELECT COUNT(*) FROM player_stats_latest)),
// 			('snapshots', now(), (SELECT COUNT(*) FROM player_stats_snapshots))
// 	`).Error

// 	if err != nil {
// 		c.log.Error("metrics collection failed", "err", err, "duration_ms", time.Since(start).Milliseconds())
// 		return
// 	}

// 	c.log.Info("metrics collection succeeded", "duration_ms", time.Since(start).Milliseconds())
// }
