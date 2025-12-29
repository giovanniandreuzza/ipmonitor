package valueobjects_test

import (
	"testing"

	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/valueobjects"
)

func TestNewIPv4Address_Valid(t *testing.T) {
	ip, err := valueobjects.NewIPv4Address("192.168.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ip.Value() != "192.168.0.1" {
		t.Fatalf("expected value to be 192.168.0.1, got %s", ip.Value())
	}
}

func TestNewIPv4Address_Empty(t *testing.T) {
	if _, err := valueobjects.NewIPv4Address(""); err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewIPv4Address_Invalid(t *testing.T) {
	if _, err := valueobjects.NewIPv4Address("invalid-ip"); err == nil {
		t.Fatal("expected error for invalid address")
	}
}
