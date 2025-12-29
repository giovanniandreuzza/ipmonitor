// Package services contains domain services.
package services

import (
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

// IPChangeDetector is a domain service for detecting IP changes.
type IPChangeDetector struct{}

// NewIPChangeDetector creates a new IP change detector.
func NewIPChangeDetector() *IPChangeDetector {
	return &IPChangeDetector{}
}

// DetectChange checks if IP has changed and generates appropriate domain event.
func (d *IPChangeDetector) DetectChange(
	stored *entities.IPMonitorState,
	current *valueobjects.IPv4Address,
) events.DomainEvent {
	if stored == nil || stored.CurrentIP() == nil {
		return events.NewIPDetectedEvent(current)
	}

	if stored.HasChanged(current) {
		return events.NewIPChangedEvent(stored.CurrentIP(), current)
	}

	return nil
}
