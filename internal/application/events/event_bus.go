// Package events provides the event bus infrastructure for the application.
package events

// Handler handles a domain event.
type Handler func(event any) error

// Bus defines event publishing behavior.
type Bus interface {
	Publish(event any) error
	Subscribe(eventName string, handler Handler)
}

// SyncEventBus is a simple in-memory synchronous event bus.
type SyncEventBus struct {
	handlers map[string][]Handler
}

// NewSyncEventBus creates a new synchronous bus.
func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{handlers: make(map[string][]Handler)}
}

// NewNoopBus returns a bus that drops all events.
func NewNoopBus() Bus {
	return noopBus{}
}

// Publish sends the event to all registered handlers.
func (b *SyncEventBus) Publish(event any) error {
	if event == nil {
		return nil
	}
	if handlers, ok := b.handlers[eventName(event)]; ok {
		for _, h := range handlers {
			if err := h(event); err != nil {
				return err
			}
		}
	}
	return nil
}

// Subscribe registers a handler for the given event name.
func (b *SyncEventBus) Subscribe(eventName string, handler Handler) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func eventName(event any) string {
	type named interface{ EventName() string }
	if e, ok := event.(named); ok {
		return e.EventName()
	}
	return "unknown"
}

type noopBus struct{}

func (noopBus) Publish(_ any) error       { return nil }
func (noopBus) Subscribe(string, Handler) {}
