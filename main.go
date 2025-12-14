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
	"albionstats/internal/tasks"

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

	ctx, cancel := signalContext(context.Background())
	defer cancel()

	// Start API server
	apiServer := api.New(db, api.Config{
		Port: cfg.APIPort,
	})

	go func() {
		addr := fmt.Sprintf(":%s", cfg.APIPort)
		log.Printf("starting API server on %s", addr)
		if err := apiServer.Run(addr); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	// Start metrics collector
	metricsCollector := tasks.NewCollector(db, tasks.CollectorConfig{
		Interval: 5 * time.Minute,
	})

	go func() {
		metricsCollector.Run(ctx)
	}()

	// Start player pollers for all regions
	regions := []string{"americas", "europe", "asia"}
	for _, region := range regions {
		playerPoller := tasks.NewPlayerPoller(db, tasks.PlayerPollerConfig{
			Region:      region,
			PageSize:    cfg.PlayerBatch,
			RatePerSec:  cfg.PlayerRate,
			UserAgent:   cfg.UserAgent,
			HTTPTimeout: cfg.HTTPTimeout,
		})

		go func(poller *tasks.PlayerPoller, regionName string) {
			log.Printf("starting player poller for region: %s", regionName)
			poller.Run(ctx)
		}(playerPoller, region)
	}

	// Start killboard pollers for all regions
	for _, region := range regions {
		kbPoller := tasks.NewKillboardPoller(db, tasks.KillboardConfig{
			PageSize:       cfg.PageSize,
			MaxPages:       cfg.MaxPages,
			EventsInterval: cfg.EventsInterval,
			Region:         region,
			HTTPTimeout:    cfg.HTTPTimeout,
			UserAgent:      cfg.UserAgent,
		})

		go func(poller *tasks.KillboardPoller, regionName string) {
			log.Printf("starting killboard poller for region: %s", regionName)
			poller.Run(ctx)
		}(kbPoller, region)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Printf("shutdown complete")
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
