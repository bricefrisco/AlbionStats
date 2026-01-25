package data_purger

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

type Purger struct {
	db       *postgres.Postgres
	interval time.Duration
	log      *slog.Logger
}

func NewPurger(db *postgres.Postgres, logger *slog.Logger, config Config) (*Purger, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if config.Interval <= 0 {
		return nil, fmt.Errorf("interval is required")
	}

	return &Purger{
		db:       db,
		interval: config.Interval,
		log:      logger.With("component", "data_purger"),
	}, nil
}

func (p *Purger) Run(ctx context.Context) {
	p.log.Info("data purger started", "interval", p.interval)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.purge(ctx)

	for {
		select {
		case <-ctx.Done():
			p.log.Info("data purger stopped")
			return
		case <-ticker.C:
			p.purge(ctx)
		}
	}
}

func (p *Purger) purge(ctx context.Context) {
	start := time.Now()
	p.log.Info("data purge started")

	if err := p.db.PurgeOldBattleData(ctx); err != nil {
		p.log.Error("data purge failed", "err", err, "duration_ms", time.Since(start).Milliseconds())
		return
	}

	p.log.Info("data purge succeeded", "duration_ms", time.Since(start).Milliseconds())
}
