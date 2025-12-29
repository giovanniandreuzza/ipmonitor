// Package file implements the IPRepository port using file persistence.
package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/entities"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/repositories"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

func getIPStoragePath() string {
	if p := os.Getenv("IP_STORAGE_PATH"); p != "" {
		return p
	}
	return "/tmp/public_ip.txt"
}

// Adapter implements IPRepository to persist IP state to file.
type Adapter struct {
	ipFilePath string
}

var _ repositories.IPRepository = (*Adapter)(nil)

// NewAdapter creates a new file persistence adapter.
func NewAdapter() *Adapter {
	return &Adapter{ipFilePath: getIPStoragePath()}
}

// Save persists the IP monitor state to a file.
func (a *Adapter) Save(state *entities.IPMonitorState) error {
	ipValue := state.CurrentIP().Value()
	dir := filepath.Dir(a.ipFilePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(a.ipFilePath, []byte(ipValue), 0o600)
}

// Load retrieves the IP monitor state from a file.
func (a *Adapter) Load() (*entities.IPMonitorState, error) {
	data, err := os.ReadFile(a.ipFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil // file not found is not an error, state is uninitialized
		}
		return nil, err
	}

	ipString := strings.TrimSpace(string(data))
	if ipString == "" {
		return nil, nil //nolint:nilnil // empty file means uninitialized state
	}

	ipv4, err := valueobjects.NewIPv4Address(ipString)
	if err != nil {
		return nil, err
	}

	return entities.NewIPMonitorState(ipv4), nil
}
