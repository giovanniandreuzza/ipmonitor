// Package ports defines the application ports (interfaces for external dependencies).
package ports

import "github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"

// NotificationPort is an application port for sending notifications about domain events.
type NotificationPort interface {
	// Notify sends a notification for a domain event
	Notify(event events.DomainEvent) error
}
