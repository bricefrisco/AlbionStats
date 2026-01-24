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
	"time"

	"albionstats/internal/api"
	"albionstats/internal/config"
	"albionstats/internal/postgres"
	"albionstats/internal/tasks"
	"albionstats/internal/tasks/battle_poller"
	"albionstats/internal/tasks/battleboard_poller"
	"albionstats/internal/tasks/killboard_poller"
	"albionstats/internal/tasks/metrics_collector"
	"albionstats/internal/tasks/player_poller"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	apiClient := tasks.NewClient()

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
		Logger:   appLogger,
	})

	go func() {
		addr := fmt.Sprintf(":%s", cfg.APIPort)
		log.Printf("starting API server on %s", addr)
		if err := server.Run(addr); err != nil {
			log.Fatalf("API server error: %v", err)
		}
	}()

	// Start metrics collector
	metricsCollector, err := metrics_collector.NewCollector(postgres, appLogger, metrics_collector.Config{
		Interval: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("metrics collector init: %v", err)
	}

	go func() {
		metricsCollector.Run(ctx)
	}()

	// Start player pollers for all regions
	regions := []string{"americas", "europe", "asia"}
	for _, region := range regions {
		playerPoller, err := player_poller.NewPlayerPoller(player_poller.Config{
			APIClient:   apiClient,
			Postgres:    postgres,
			Logger:      appLogger,
			Region:      region,
			BatchSize:   cfg.PlayerBatch,
			WorkerCount: cfg.PlayerWorkerCount,
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
			Region:         region,
			APIClient:      apiClient,
			Postgres:       postgres,
			Logger:         appLogger,
			PageSize:       cfg.EventsPageSize,
			MaxPages:       cfg.EventsMaxPages,
			EventsInterval: cfg.EventsInterval,
		})
		if err != nil {
			log.Fatalf("killboard poller init (%s): %v", region, err)
		}

		go func(poller *killboard_poller.KillboardPoller, regionName string) {
			log.Printf("starting killboard poller for region: %s", regionName)
			poller.Run()
		}(kbPoller, region)
	}

	// Start battleboard poller for all regions
	for _, region := range regions {
		battleboardPoller := battleboard_poller.NewBattleboardPoller(battleboard_poller.Config{
			APIClient:      apiClient,
			Postgres:       postgres,
			Logger:         appLogger,
			Region:         region,
			PageSize:       cfg.BattleboardPageSize,
			MaxPages:       cfg.BattleboardMaxPages,
			EventsInterval: cfg.BattleboardInterval,
		})

		go func(poller *battleboard_poller.BattleboardPoller, regionName string) {
			log.Printf("starting battleboard poller for region: %s", regionName)
			poller.Run()
		}(battleboardPoller, region)
	}

	// Start battle poller for all regions
	for _, region := range regions {
		battlePoller := battle_poller.NewBattlePoller(battle_poller.Config{
			Region:    region,
			APIClient: apiClient,
			Postgres:  postgres,
			Logger:    appLogger,
		})

		go func(poller *battle_poller.BattlePoller, regionName string) {
			log.Printf("starting battle poller for region: %s", regionName)
			poller.Run()
		}(battlePoller, region)
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
