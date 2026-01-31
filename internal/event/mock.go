package event

import (
	"context"
	"sync"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	MockPublisher simulates Kafka: stores events in memory for testing and observability
	In production, replace with a real Kafka producer
*/

type MockPublisher struct {
	mu     sync.Mutex
	events []*domain.DomainEvent
}

// Creates a new in-memory event publisher
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{events: make([]*domain.DomainEvent, 0)}
}

// Appends the event to in-memory storage (simulates Kafka send)
func (p *MockPublisher) Publish(ctx context.Context, event *domain.DomainEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, event)
	return nil
}

// Closes the mock publisher
func (p *MockPublisher) Close() error { return nil }

// Returns a copy of all published events (for tests and debugging)
func (p *MockPublisher) Events() []*domain.DomainEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]*domain.DomainEvent, len(p.events))
	copy(out, p.events)
	return out
}
