package event

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

// Publishes domain events (Kafka can be added here)
type Publisher interface {
	Publish(ctx context.Context, event *domain.DomainEvent) error
	Close() error
}
