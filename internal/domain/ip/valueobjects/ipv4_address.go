// Package valueobjects contains domain value objects.
package valueobjects

import (
	"net"

	domainerrors "github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/errors"
)

// IPv4Address is a value object representing a valid IPv4 address.
type IPv4Address struct {
	value string
}

// NewIPv4Address creates a new IPv4Address with validation.
func NewIPv4Address(address string) (*IPv4Address, error) {
	if address == "" {
		return nil, domainerrors.EmptyIPv4Error()
	}

	ip := net.ParseIP(address)
	if ip == nil || ip.To4() == nil {
		return nil, domainerrors.InvalidIPv4Error(address)
	}

	return &IPv4Address{value: address}, nil
}

// Value returns the string representation of the IPv4 address.
func (ip *IPv4Address) Value() string {
	return ip.value
}

// Equals checks if two IPv4 addresses are equal.
func (ip *IPv4Address) Equals(other *IPv4Address) bool {
	if other == nil {
		return false
	}
	return ip.value == other.value
}

// String implements the Stringer interface.
func (ip *IPv4Address) String() string {
	return ip.value
}
