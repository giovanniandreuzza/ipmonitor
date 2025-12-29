// Package dto contains data transfer objects for the monitor IP use case.
package dto

import "github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"

// MonitorIPCommand represents the input for monitoring IP use case.
type MonitorIPCommand struct {
	// Currently no input needed, but structured for future extension
}

// Validate validates the input command. Future fields go here.
func (c *MonitorIPCommand) Validate() error {
	return nil
}

// MonitorIPResult represents the output of monitoring IP use case.
type MonitorIPResult struct {
	IPChanged   bool
	OldIP       *valueobjects.IPv4Address
	NewIP       *valueobjects.IPv4Address
	WasNotified bool
	Error       error
}

// NewMonitorIPResult creates a result for when IP hasn't changed.
func NewMonitorIPResult(changed bool, oldIP, newIP *valueobjects.IPv4Address, notified bool) *MonitorIPResult {
	return &MonitorIPResult{
		IPChanged:   changed,
		OldIP:       oldIP,
		NewIP:       newIP,
		WasNotified: notified,
	}
}

// NewMonitorIPError creates a result for errors.
func NewMonitorIPError(err error) *MonitorIPResult {
	return &MonitorIPResult{
		Error: err,
	}
}
