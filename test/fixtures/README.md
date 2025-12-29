# Test Fixtures

This directory contains test data and fixtures for integration testing.

## Structure

```
fixtures/
├── valid_ips.txt       # Valid IP addresses for testing
├── invalid_ips.txt     # Invalid IP addresses for testing
└── mock_responses/     # Mock API responses
```

## Usage

Test files can reference fixtures using relative paths from the test files.
