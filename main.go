package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	// apiServer := api.New(db, api.Config{
	// 	Port: cfg.APIPort,
	// })

	// go func() {
	// 	addr := fmt.Sprintf(":%s", cfg.APIPort)
	// 	log.Printf("starting API server on %s", addr)
	// 	if err := apiServer.Run(addr); err != nil {
	// 		log.Printf("API server error: %v", err)
	// 	}
	// }()

	// Start metrics collector
	// metricsCollector := tasks.NewCollector(db, tasks.CollectorConfig{
	// 	Interval: 5 * time.Minute,
	// })

	// go func() {
	// 	metricsCollector.Run(ctx)
	// }()

	// Start player poller
	// playerPoller := tasks.NewPlayerPoller(db, tasks.PlayerPollerConfig{
	// 	APIBase:     cfg.APIBase,
	// 	PageSize:    cfg.PlayerBatch,
	// 	RatePerSec:  cfg.PlayerRate,
	// 	UserAgent:   cfg.UserAgent,
	// 	HTTPTimeout: cfg.HTTPTimeout,
	// })

	// go func() {
	// 	playerPoller.Run(ctx)
	// }()

	// Start killboard poller
	kbPoller := tasks.NewKillboardPoller(db, tasks.KillboardConfig{
		APIBase:        cfg.APIBase,
		PageSize:       cfg.PageSize,
		MaxPages:       cfg.MaxPages,
		EventsInterval: cfg.EventsInterval,
		Region:         cfg.Region,
		HTTPTimeout:    cfg.HTTPTimeout,
		UserAgent:      cfg.UserAgent,
	})

	go func() {
		kbPoller.Run(ctx)
	}()

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
