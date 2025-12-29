// Package ipify implements the PublicIPProvider port using the ipify API.
package ipify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/giovanniandreuzza/ipmonitor/internal/application/ports"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

const (
	ipifyAPIURL     = "https://api.ipify.org"
	httpTimeoutSecs = 10
)

// Adapter implements PublicIPProvider to fetch IP from ipify API.
type Adapter struct {
	apiURL string
}

var _ ports.PublicIPProvider = (*Adapter)(nil)

// NewAdapter creates a new ipify adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		apiURL: ipifyAPIURL,
	}
}

// GetPublicIP retrieves the current public IPv4 address from ipify API.
func (a *Adapter) GetPublicIP() (*valueobjects.IPv4Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutSecs*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public IP: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ipString := strings.TrimSpace(string(body))
	ip, err := valueobjects.NewIPv4Address(ipString)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address received: %w", err)
	}

	return ip, nil
}
