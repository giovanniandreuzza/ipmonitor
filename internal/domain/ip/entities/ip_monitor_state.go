// Package entities contains domain entities.
package entities

import (
	"time"

	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

// IPMonitorState represents the current state of IP monitoring.
type IPMonitorState struct {
	currentIP   *valueobjects.IPv4Address
	lastChecked time.Time
}

// NewIPMonitorState creates a new IP monitor state.
func NewIPMonitorState(ip *valueobjects.IPv4Address) *IPMonitorState {
	return &IPMonitorState{
		currentIP:   ip,
		lastChecked: time.Now(),
	}
}

// CurrentIP returns the current IP address.
func (s *IPMonitorState) CurrentIP() *valueobjects.IPv4Address {
	return s.currentIP
}

// LastChecked returns when the IP was last checked.
func (s *IPMonitorState) LastChecked() time.Time {
	return s.lastChecked
}

// UpdateIP updates the current IP and timestamp.
func (s *IPMonitorState) UpdateIP(newIP *valueobjects.IPv4Address) {
	s.currentIP = newIP
	s.lastChecked = time.Now()
}

// HasChanged checks if the IP has changed from the stored state.
func (s *IPMonitorState) HasChanged(newIP *valueobjects.IPv4Address) bool {
	if s.currentIP == nil {
		return true
	}
	return !s.currentIP.Equals(newIP)
}
