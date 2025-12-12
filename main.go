package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albionstats/internal/api"
	"albionstats/internal/config"
	"albionstats/internal/killboard"
	"albionstats/internal/metrics"
	"albionstats/internal/playerpoller"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// Start API server
	apiServer := api.New(db, api.Config{
		Port: cfg.APIPort,
	})

	go func() {
		addr := fmt.Sprintf(":%s", cfg.APIPort)
		log.Printf("starting API server on %s", addr)
		if err := apiServer.Start(addr); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	kbPoller := killboard.New(db, killboard.Config{
		APIBase:        cfg.APIBase,
		PageSize:       cfg.PageSize,
		MaxPages:       cfg.MaxPages,
		EventsInterval: cfg.EventsInterval,
		Region:         cfg.Region,
		HTTPTimeout:    cfg.HTTPTimeout,
		UserAgent:      cfg.UserAgent,
	})

	playerPoller := playerpoller.New(db, playerpoller.Config{
		APIBase:     cfg.APIBase,
		PageSize:    cfg.PlayerBatch,
		RatePerSec:  cfg.PlayerRate,
		UserAgent:   cfg.UserAgent,
		HTTPTimeout: cfg.HTTPTimeout,
	})

	ctx, cancel := signalContext(context.Background())
	defer cancel()

	// Start metrics collector
	metricsCollector := metrics.New(db, metrics.Config{
		Interval: 5 * time.Minute,
	})

	go func() {
		metricsCollector.Start(ctx)
	}()

	// Player poller runs continuously; it rate-limits internally.
	go func() {
		for {
			if err := playerPoller.Run(ctx); err != nil {
				log.Printf("player poller error: %v", err)
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	ticker := time.NewTicker(cfg.EventsInterval)
	defer ticker.Stop()

	// Run killboard poller immediately once
	log.Printf("loop: running killboard poller")
	if err := kbPoller.Run(ctx); err != nil {
		log.Printf("poller error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		log.Printf("loop: running killboard poller")
		if err := kbPoller.Run(ctx); err != nil {
			log.Printf("poller error: %v", err)
		}
	}
}

func signalContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
