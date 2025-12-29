// Package telegram implements the NotificationPort using the Telegram API.
package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/giovanniandreuzza/ipmonitor/internal/application/ports"
	"github.com/giovanniandreuzza/ipmonitor/internal/domain/ip/events"
)

const httpTimeoutSecs = 10

// Adapter implements NotificationPort to send Telegram messages.
type Adapter struct {
	botToken string
	chatID   string
}

var _ ports.NotificationPort = (*Adapter)(nil)

// NewAdapter creates a new Telegram notification adapter.
func NewAdapter(botToken, chatID string) *Adapter {
	return &Adapter{
		botToken: botToken,
		chatID:   chatID,
	}
}

// Notify sends a notification based on the domain event.
func (a *Adapter) Notify(event events.DomainEvent) error {
	message := a.formatMessage(event)
	return a.sendMessage(message)
}

func (a *Adapter) formatMessage(event events.DomainEvent) string {
	switch e := event.(type) {
	case *events.IPChangedEvent:
		return fmt.Sprintf(
			"🌐 *IP Monitor Alert*\n\n🔄 Your public IP has changed:\n\nOld IP: %s\nNew IP: %s",
			e.OldIP().String(),
			e.NewIP().String(),
		)
	case *events.IPDetectedEvent:
		return fmt.Sprintf(
			"🌐 *IP Monitor Alert*\n\n✅ Initial IP detected: %s",
			e.IP().String(),
		)
	default:
		return fmt.Sprintf("IP Monitor: %s", event.EventName())
	}
}

func (a *Adapter) sendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutSecs*time.Second)
	defer cancel()

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.botToken)

	data := url.Values{}
	data.Set("chat_id", a.chatID)
	data.Set("text", message)
	data.Set("parse_mode", "Markdown")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
