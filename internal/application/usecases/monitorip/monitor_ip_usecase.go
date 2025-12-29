// Package monitorip implements the monitor IP use case.
package monitorip

import (
	"log/slog"

	appevents "github.com/giovanniandreuzza/ipmonitor/internal/application/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/ports"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip/dto"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/repositories"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/services"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

// UseCase encapsulates the business logic for monitoring IP changes.
type UseCase struct {
	ipProvider     ports.PublicIPProvider
	ipRepository   repositories.IPRepository
	notifier       ports.NotificationPort
	changeDetector *services.IPChangeDetector
	eventBus       appevents.Bus
}

// NewUseCase creates a new monitor IP use case.
func NewUseCase(
	ipProvider ports.PublicIPProvider,
	ipRepository repositories.IPRepository,
	notifier ports.NotificationPort,
	eventBus appevents.Bus,
) *UseCase {
	if eventBus == nil {
		eventBus = appevents.NewNoopBus()
	}
	return &UseCase{
		ipProvider:     ipProvider,
		ipRepository:   ipRepository,
		notifier:       notifier,
		changeDetector: services.NewIPChangeDetector(),
		eventBus:       eventBus,
	}
}

// Execute runs the IP monitoring use case.
func (uc *UseCase) Execute(cmd *dto.MonitorIPCommand) *dto.MonitorIPResult {
	if err := cmd.Validate(); err != nil {
		return dto.NewMonitorIPError(err)
	}

	// 1. Get current public IP
	currentIP, err := uc.ipProvider.GetPublicIP()
	if err != nil {
		return dto.NewMonitorIPError(err)
	}

	// 2. Load stored state
	storedState, err := uc.ipRepository.Load()
	if err != nil {
		return dto.NewMonitorIPError(err)
	}

	// 3. Detect changes using domain service
	event := uc.changeDetector.DetectChange(storedState, currentIP)

	// 4. If no change, return early
	if event == nil {
		slog.Info("IP has not changed") //nolint:sloglint // logger set by application initialization
		return dto.NewMonitorIPResult(false, currentIP, currentIP, false)
	}

	// 5. Handle domain event
	var oldIP *entities.IPMonitorState
	if storedState != nil {
		oldIP = storedState
	}

	// 6. Send notification
	if notifyErr := uc.notifier.Notify(event); notifyErr != nil {
		return dto.NewMonitorIPError(notifyErr)
	}

	if pubErr := uc.eventBus.Publish(event); pubErr != nil {
		return dto.NewMonitorIPError(pubErr)
	}

	//nolint:sloglint // logger set by application initialization
	slog.Info(
		"IP change detected",
		slog.String("event", event.EventName()),
	)

	// 7. Update repository
	newState := entities.NewIPMonitorState(currentIP)
	if saveErr := uc.ipRepository.Save(newState); saveErr != nil {
		return dto.NewMonitorIPError(saveErr)
	}

	// 8. Build result
	var oldIPValue *valueobjects.IPv4Address
	if oldIP != nil {
		oldIPValue = oldIP.CurrentIP()
	}

	switch e := event.(type) {
	case *events.IPChangedEvent:
		return dto.NewMonitorIPResult(true, e.OldIP(), e.NewIP(), true)
	case *events.IPDetectedEvent:
		return dto.NewMonitorIPResult(true, oldIPValue, e.IP(), true)
	}

	return dto.NewMonitorIPResult(false, currentIP, currentIP, false)
}
