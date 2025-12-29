// Package events contains domain events.
package events

import (
	"time"

	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	OccurredAt() time.Time
	EventName() string
}

// IPChangedEvent represents a domain event when IP address changes.
type IPChangedEvent struct {
	oldIP      *valueobjects.IPv4Address
	newIP      *valueobjects.IPv4Address
	occurredAt time.Time
}

// NewIPChangedEvent creates a new IP changed event.
func NewIPChangedEvent(oldIP, newIP *valueobjects.IPv4Address) *IPChangedEvent {
	return &IPChangedEvent{
		oldIP:      oldIP,
		newIP:      newIP,
		occurredAt: time.Now(),
	}
}

// OldIP returns the previous IP address.
func (e *IPChangedEvent) OldIP() *valueobjects.IPv4Address {
	return e.oldIP
}

// NewIP returns the new IP address.
func (e *IPChangedEvent) NewIP() *valueobjects.IPv4Address {
	return e.newIP
}

// OccurredAt returns when the event occurred.
func (e *IPChangedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventName returns the name of the event.
func (e *IPChangedEvent) EventName() string {
	return "IPChanged"
}

// IPDetectedEvent represents the initial IP detection event.
type IPDetectedEvent struct {
	ip         *valueobjects.IPv4Address
	occurredAt time.Time
}

// NewIPDetectedEvent creates a new IP detected event.
func NewIPDetectedEvent(ip *valueobjects.IPv4Address) *IPDetectedEvent {
	return &IPDetectedEvent{
		ip:         ip,
		occurredAt: time.Now(),
	}
}

// IP returns the detected IP address.
func (e *IPDetectedEvent) IP() *valueobjects.IPv4Address {
	return e.ip
}

// OccurredAt returns when the event occurred.
func (e *IPDetectedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventName returns the name of the event.
func (e *IPDetectedEvent) EventName() string {
	return "IPDetected"
}
