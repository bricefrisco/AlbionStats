package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"albionstats/internal/api"
	"albionstats/internal/config"
	"albionstats/internal/postgres"
	"albionstats/internal/tasks/killboard_poller"
	"albionstats/internal/tasks/player_poller"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	apiClient := api.NewClient()

	postgres, err := postgres.NewPostgresDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("postgres database: %v", err)
	}

	appLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // Remove timestamp since syslog provides it
			}
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				// Extract just the filename from the full path
				if lastSlash := strings.LastIndex(source.File, "/"); lastSlash >= 0 {
					source.File = source.File[lastSlash+1:]
				}
				return a
			}
			return a
		},
	}))

	ctx, cancel := signalContext(context.Background())
	defer cancel()

	server := api.NewServer(api.Config{
		Postgres: postgres,
	})

	go func() {
		addr := fmt.Sprintf(":%s", cfg.APIPort)
		log.Printf("starting API server on %s", addr)
		if err := server.Run(addr); err != nil {
			log.Fatalf("API server error: %v", err)
		}
	}()

	// Start metrics collector
	// metricsCollector, err := tasks.NewCollector(db, appLogger, tasks.CollectorConfig{
	// 	Interval: 5 * time.Minute,
	// })
	// if err != nil {
	// 	log.Fatalf("metrics collector init: %v", err)
	// }

	// go func() {
	// 	metricsCollector.Run(ctx)
	// }()

	// Start player pollers for all regions
	regions := []string{"americas", "europe", "asia"}
	for _, region := range regions {
		playerPoller, err := player_poller.NewPlayerPoller(player_poller.Config{
			APIClient:  apiClient,
			Postgres:   postgres,
			Logger:     appLogger,
			Region:     region,
			BatchSize:  cfg.PlayerBatch,
			RatePerSec: cfg.PlayerRate,
		})
		if err != nil {
			log.Fatalf("player poller init (%s): %v", region, err)
		}

		go func(poller *player_poller.PlayerPoller, regionName string) {
			log.Printf("starting player poller for region: %s", regionName)
			poller.Run()
		}(playerPoller, region)
	}

	// Start killboard pollers for all regions
	for _, region := range regions {
		kbPoller, err := killboard_poller.NewKillboardPoller(killboard_poller.Config{
			PageSize:       cfg.PageSize,
			MaxPages:       cfg.MaxPages,
			EventsInterval: cfg.EventsInterval,
			Region:         region,
			APIClient:      apiClient,
			Postgres:       postgres,
			Logger:         appLogger,
		})
		if err != nil {
			log.Fatalf("killboard poller init (%s): %v", region, err)
		}

		go func(poller *killboard_poller.KillboardPoller, regionName string) {
			log.Printf("starting killboard poller for region: %s", regionName)
			poller.Run()
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
