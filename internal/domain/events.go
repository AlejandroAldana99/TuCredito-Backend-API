package domain

import (
	"encoding/json"
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

// Payload for the CreditCreated event
type CreditCreatedPayload struct {
	CreditID   string       `json:"credit_id"`
	ClientID   string       `json:"client_id"`
	BankID     string       `json:"bank_id"`
	CreditType CreditType   `json:"credit_type"`
	Status     CreditStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
}

// Payload for the CreditApproved event
type CreditApprovedPayload struct {
	CreditID   string    `json:"credit_id"`
	ClientID   string    `json:"client_id"`
	BankID     string    `json:"bank_id"`
	ApprovedAt time.Time `json:"approved_at"`
}

// Payload for the CreditRejected event
type CreditRejectedPayload struct {
	CreditID   string    `json:"credit_id"`
	ClientID   string    `json:"client_id"`
	BankID     string    `json:"bank_id"`
	RejectedAt time.Time `json:"rejected_at"`
}

// Serializes a payload to JSON for the event envelope
func MarshalPayload(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
