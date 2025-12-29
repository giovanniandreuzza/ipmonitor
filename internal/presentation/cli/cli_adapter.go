// Package cli provides the CLI adapter for the IP monitor application.
package cli

import (
	"log/slog"

	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip/dto"
)

// Adapter is a primary adapter for CLI interaction.
type Adapter struct {
	monitorIPUseCase *monitorip.UseCase
}

// NewAdapter creates a new CLI adapter.
func NewAdapter(monitorIPUseCase *monitorip.UseCase) *Adapter {
	return &Adapter{
		monitorIPUseCase: monitorIPUseCase,
	}
}

// Run executes the IP monitoring through CLI.
func (a *Adapter) Run() error {
	cmd := &dto.MonitorIPCommand{}
	result := a.monitorIPUseCase.Execute(cmd)

	if result.Error != nil {
		//nolint:sloglint // logger set by application initialization
		slog.Error(
			"monitor IP failed",
			slog.Any("error", result.Error),
		)
		return result.Error
	}

	if result.IPChanged {
		slog.Info("IP changed, notification sent") //nolint:sloglint // logger set by application initialization
	} else {
		slog.Info("IP has not changed") //nolint:sloglint // logger set by application initialization
	}

	return nil
}
