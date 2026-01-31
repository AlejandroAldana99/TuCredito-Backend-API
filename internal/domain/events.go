package domain

import (
	"time"
)

// Type of domain events for event-driven architecture (CreditCreated, CreditApproved, CreditRejected)
type EventType string

const (
	EventCreditCreated  EventType = "CreditCreated"
	EventCreditApproved EventType = "CreditApproved"
	EventCreditRejected EventType = "CreditRejected"
)

// Envelope for all domain events emitted by the service
type DomainEvent struct {
	ID         string    `json:"id"`
	Type       EventType `json:"type"`
	Payload    []byte    `json:"payload"`
	OccurredAt time.Time `json:"occurred_at"`
}
