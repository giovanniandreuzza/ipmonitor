package errors //nolint:revive // intentional package name conflicts with stdlib for clarity

import "fmt"

// DomainError provides typed domain errors with codes.
type DomainError interface {
	error
	Code() string
}

type baseError struct {
	code string
	msg  string
}

func (e *baseError) Error() string { return e.msg }
func (e *baseError) Code() string  { return e.code }

// InvalidIPv4Error indicates an invalid IPv4 address value.
func InvalidIPv4Error(value string) DomainError {
	return &baseError{code: "INVALID_IPV4", msg: fmt.Sprintf("invalid IPv4 address: %s", value)}
}

// EmptyIPv4Error indicates an empty IPv4 value.
func EmptyIPv4Error() DomainError {
	return &baseError{code: "EMPTY_IPV4", msg: "IPv4 address cannot be empty"}
}

// RepositoryLoadError indicates persistence retrieval issues.
func RepositoryLoadError(err error) DomainError {
	return &baseError{code: "REPO_LOAD_FAILED", msg: fmt.Sprintf("failed to load IP state: %v", err)}
}

// NotificationError indicates notification failures.
func NotificationError(err error) DomainError {
	return &baseError{code: "NOTIFICATION_FAILED", msg: fmt.Sprintf("failed to send notification: %v", err)}
}
