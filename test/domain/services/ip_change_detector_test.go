package services_test

import (
	"testing"

	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/services"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

func TestDetectChange_FirstDetection(t *testing.T) {
	detector := services.NewIPChangeDetector()
	current, _ := valueobjects.NewIPv4Address("1.1.1.1")

	event := detector.DetectChange(nil, current)
	if _, ok := event.(*events.IPDetectedEvent); !ok {
		t.Fatalf("expected IPDetectedEvent, got %T", event)
	}
}

func TestDetectChange_Changed(t *testing.T) {
	detector := services.NewIPChangeDetector()
	storedIP, _ := valueobjects.NewIPv4Address("1.1.1.1")
	currentIP, _ := valueobjects.NewIPv4Address("2.2.2.2")
	stored := entities.NewIPMonitorState(storedIP)

	event := detector.DetectChange(stored, currentIP)
	if _, ok := event.(*events.IPChangedEvent); !ok {
		t.Fatalf("expected IPChangedEvent, got %T", event)
	}
}

func TestDetectChange_NoChange(t *testing.T) {
	detector := services.NewIPChangeDetector()
	storedIP, _ := valueobjects.NewIPv4Address("1.1.1.1")
	currentIP, _ := valueobjects.NewIPv4Address("1.1.1.1")
	stored := entities.NewIPMonitorState(storedIP)

	event := detector.DetectChange(stored, currentIP)
	if event != nil {
		t.Fatalf("expected no event, got %T", event)
	}
}
