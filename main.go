package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albionstats/internal/config"
	"albionstats/internal/killboard"

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

	client := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}

	poller := killboard.New(client, db, killboard.Config{
		APIBase:        cfg.APIBase,
		PageSize:       cfg.PageSize,
		MaxPages:       cfg.MaxPages,
		EventsInterval: cfg.EventsInterval,
		Region:         cfg.Region,
		UserAgent:      cfg.UserAgent,
	})

	ctx, cancel := signalContext(context.Background())
	defer cancel()

	ticker := time.NewTicker(cfg.EventsInterval)
	defer ticker.Stop()

	for {
		log.Printf("loop: running poller")
		if err := poller.Run(ctx); err != nil {
			log.Printf("poller error: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
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
