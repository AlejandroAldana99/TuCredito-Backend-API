package decision

import (
	"context"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	Rule defines a single routing or eligibility rule
	Extensible for ad-tech
*/

type Rule interface {
	Evaluate(ctx context.Context, input *EligibilityInput) (approved bool, priority int, score float64)
	Name() string
}

// Input for the decision engine
type EligibilityInput struct {
	Client     *domain.Client
	Bank       *domain.Bank
	MinPayment float64
	MaxPayment float64
	TermMonths int
	CreditType domain.CreditType
}

// Result of the eligibility evaluation
type EligibilityResult struct {
	Approved bool    `json:"approved"`
	Priority int     `json:"priority"`
	Score    float64 `json:"score"`
	RuleName string  `json:"rule_name"`
}

// Routes credit applications and applies config rules (waterfall/priority)
type Engine interface {
	Evaluate(ctx context.Context, input *EligibilityInput) (*EligibilityResult, error)
	RegisterRule(rule Rule)
}
