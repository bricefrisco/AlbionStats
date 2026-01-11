package metrics_collector

import (
	"albionstats/internal/postgres"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Config struct {
	Interval time.Duration
}

type Collector struct {
	db     *postgres.Postgres
	config Config
	log    *slog.Logger
}

func NewCollector(db *postgres.Postgres, logger *slog.Logger, config Config) (*Collector, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &Collector{
		db:     db,
		config: config,
		log:    logger.With("component", "metrics_collector"),
	}, nil
}

func (c *Collector) Run(ctx context.Context) {
	c.log.Info("collector started", "interval", c.config.Interval)

	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	// Run once immediately
	c.collect(ctx)

	for {
		select {
		case <-ctx.Done():
			c.log.Info("collector stopped")
			return
		case <-ticker.C:
			c.collect(ctx)
		}
	}
}

func (c *Collector) collect(ctx context.Context) {
	start := time.Now()
	c.log.Info("metrics collection started")

	if err := c.db.InsertMetrics(ctx); err != nil {
		c.log.Error("metrics collection failed", "err", err, "duration_ms", time.Since(start).Milliseconds())
		return
	}

	c.log.Info("metrics collection succeeded", "duration_ms", time.Since(start).Milliseconds())
}
