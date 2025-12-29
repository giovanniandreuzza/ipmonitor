package application_test

import (
	"errors"
	"testing"

	appevents "github.com/giovanniandreuzza/ipmonitor/internal/application/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip"
	"github.com/giovanniandreuzza/ipmonitor/internal/application/usecases/monitorip/dto"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

func TestUseCase_IPChangedFlow(t *testing.T) {
	ipProvider := &mockIPProvider{ip: "2.2.2.2"}
	repo := &mockRepo{state: entities.NewIPMonitorState(mustIP("1.1.1.1"))}
	notifier := &mockNotifier{}
	eventBus := appevents.NewSyncEventBus()

	uc := monitorip.NewUseCase(ipProvider, repo, notifier, eventBus)
	result := uc.Execute(&dto.MonitorIPCommand{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if !result.IPChanged || !result.WasNotified {
		t.Fatalf("expected IP change and notification")
	}
	if repo.saved == nil || repo.saved.CurrentIP().Value() != "2.2.2.2" {
		t.Fatalf("repository should have new IP saved")
	}
}

func TestUseCase_NoChange(t *testing.T) {
	ipProvider := &mockIPProvider{ip: "1.1.1.1"}
	repo := &mockRepo{state: entities.NewIPMonitorState(mustIP("1.1.1.1"))}
	notifier := &mockNotifier{}
	eventBus := appevents.NewSyncEventBus()

	uc := monitorip.NewUseCase(ipProvider, repo, notifier, eventBus)
	result := uc.Execute(&dto.MonitorIPCommand{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if result.IPChanged {
		t.Fatalf("expected no IP change")
	}
	if notifier.called {
		t.Fatalf("notifier should not be called when no change")
	}
}

// --- Mocks ---

type mockIPProvider struct{ ip string }

func (m *mockIPProvider) GetPublicIP() (*valueobjects.IPv4Address, error) {
	return mustIP(m.ip), nil
}

type mockRepo struct {
	state *entities.IPMonitorState
	saved *entities.IPMonitorState
}

func (m *mockRepo) Save(state *entities.IPMonitorState) error {
	m.saved = state
	return nil
}

func (m *mockRepo) Load() (*entities.IPMonitorState, error) {
	return m.state, nil
}

type mockNotifier struct{ called bool }

func (m *mockNotifier) Notify(event events.DomainEvent) error {
	m.called = true
	if event == nil {
		return errors.New("event required")
	}
	return nil
}

func mustIP(value string) *valueobjects.IPv4Address {
	ip, err := valueobjects.NewIPv4Address(value)
	if err != nil {
		panic(err)
	}
	return ip
}
