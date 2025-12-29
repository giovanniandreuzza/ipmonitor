// Package main is the entry point for the IP monitor application.
package main

import (
	"log/slog"
	"os"

	appevents "github.com/giovanniandreuzza/ipmonitor/internal/application/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip"
	"github.com/giovanniandreuzza/ipmonitor/internal/infrastructure/adapters/http/ipify"
	"github.com/giovanniandreuzza/ipmonitor/internal/infrastructure/adapters/notification/telegram"
	"github.com/giovanniandreuzza/ipmonitor/internal/infrastructure/adapters/persistence/file"
	"github.com/giovanniandreuzza/ipmonitor/internal/infrastructure/config"
	"github.com/giovanniandreuzza/ipmonitor/internal/infrastructure/logging"
	"github.com/giovanniandreuzza/ipmonitor/internal/presentation/cli"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		//nolint:sloglint // global logger set in main
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger := logging.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	// Initialize adapters (infrastructure implementations)
	ipProvider := ipify.NewAdapter()
	ipRepository := file.NewAdapter()
	notifier := telegram.NewAdapter(cfg.TelegramBotToken, cfg.TelegramChatID)
	eventBus := appevents.NewSyncEventBus()

	// Initialize application use case
	monitorIPUseCase := monitorip.NewUseCase(
		ipProvider,
		ipRepository,
		notifier,
		eventBus,
	)

	// Initialize primary adapter (CLI - driving/inbound)
	cliAdapter := cli.NewAdapter(monitorIPUseCase)

	// Run the application
	if runErr := cliAdapter.Run(); runErr != nil {
		//nolint:sloglint // global logger set in main
		slog.Error(
			"Error running IP monitor",
			slog.Any("error", runErr),
		)
		os.Exit(1)
	}
}
