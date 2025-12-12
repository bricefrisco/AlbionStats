package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albionstats/internal/config"
	"albionstats/internal/killboard"
	"albionstats/internal/playerpoller"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

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
