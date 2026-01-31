package decision

import (
	"context"
	"sync"

	"github.com/tucredito/backend-api/internal/domain"
)

/*
	RuleEngine implements Engine with config rules
	Waterfall/priority order matters for the evaluation
*/

type RuleEngine struct {
	mu    sync.RWMutex
	rules []Rule
}

// Creates a new decision engine
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{rules: make([]Rule, 0)}
}

// Adds a rule to the engine (order matters for waterfall)
func (e *RuleEngine) RegisterRule(rule Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, rule)
}

// Runs rules in order and returns the first approval or the last result
func (e *RuleEngine) Evaluate(ctx context.Context, input *EligibilityInput) (*EligibilityResult, error) {
	e.mu.RLock()
	rules := make([]Rule, len(e.rules))
	copy(rules, e.rules)
	e.mu.RUnlock()

	if len(rules) == 0 {
		return &EligibilityResult{Approved: false, Priority: 0, Score: 0, RuleName: "none"}, nil
	}

	var last *EligibilityResult
	for _, rule := range rules {
		approved, priority, score := rule.Evaluate(ctx, input)
		last = &EligibilityResult{
			Approved: approved,
			Priority: priority,
			Score:    score,
			RuleName: rule.Name(),
		}
		if approved {
			return last, nil
		}
	}

	return last, nil
}

/*
	PaymentRangeRule approves when payment is within min/max and term is positive
	MinPayment <= 0 || MaxPayment < input.MinPayment || TermMonths <= 0 => false, 10, 0
	MaxPayment > 0 => score = MinPayment / MaxPayment
	Otherwise, score = 1.0
*/

type PaymentRangeRule struct{}

func (PaymentRangeRule) Name() string { return "PaymentRangeRule" }

func (PaymentRangeRule) Evaluate(ctx context.Context, input *EligibilityInput) (bool, int, float64) {
	if input == nil {
		return false, 0, 0
	}

	if input.MinPayment <= 0 || input.MaxPayment < input.MinPayment || input.TermMonths <= 0 {
		return false, 10, 0
	}

	score := 1.0
	if input.MaxPayment > 0 {
		score = input.MinPayment / input.MaxPayment
	}

	return true, 10, score
}

/*
	BankTypeRule gives priority to GOVERNMENT banks
	Bank.Type == domain.BankTypeGovernment => true, 8, 0.8
	Otherwise, true, 5, 0.5
*/

type BankTypeRule struct{}

func (BankTypeRule) Name() string { return "BankTypeRule" }

func (BankTypeRule) Evaluate(ctx context.Context, input *EligibilityInput) (bool, int, float64) {
	if input == nil || input.Bank == nil {
		return false, 0, 0
	}

	priority := 5
	score := 0.5
	if input.Bank.Type == domain.BankTypeGovernment {
		priority = 8
		score = 0.8
	}

	return true, priority, score
}
