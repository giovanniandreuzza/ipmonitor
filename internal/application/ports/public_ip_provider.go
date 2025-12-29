package ports

import "github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"

// PublicIPProvider is a port for fetching the public IP address.
type PublicIPProvider interface {
	// GetPublicIP retrieves the current public IPv4 address
	GetPublicIP() (*valueobjects.IPv4Address, error)
}
