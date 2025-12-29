// Package repositories defines repository interfaces.
package repositories

import "github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"

// IPRepository is a repository interface for persisting IP state.
type IPRepository interface {
	// Save persists the current IP monitor state
	Save(state *entities.IPMonitorState) error

	// Load retrieves the stored IP monitor state
	Load() (*entities.IPMonitorState, error)
}
