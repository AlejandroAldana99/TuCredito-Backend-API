package decision

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tucredito/backend-api/internal/domain"
)

func TestRuleEngine_Evaluate_EmptyRules(t *testing.T) {
	engine := NewRuleEngine()
	result, err := engine.Evaluate(context.Background(), &EligibilityInput{
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 12,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Approved)
	assert.Equal(t, "none", result.RuleName)
}

func TestRuleEngine_Evaluate_PaymentRangeRule(t *testing.T) {
	engine := NewRuleEngine()
	engine.RegisterRule(PaymentRangeRule{})
	result, err := engine.Evaluate(context.Background(), &EligibilityInput{
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 12,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Approved)
	assert.Equal(t, "PaymentRangeRule", result.RuleName)
	assert.GreaterOrEqual(t, result.Score, 0.0)
}

func TestRuleEngine_Evaluate_PaymentRangeRule_Invalid(t *testing.T) {
	engine := NewRuleEngine()
	engine.RegisterRule(PaymentRangeRule{})
	result, err := engine.Evaluate(context.Background(), &EligibilityInput{
		MinPayment: 500,
		MaxPayment: 100,
		TermMonths: 12,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Approved)
}

func TestRuleEngine_Evaluate_Waterfall(t *testing.T) {
	engine := NewRuleEngine()
	engine.RegisterRule(PaymentRangeRule{})
	engine.RegisterRule(BankTypeRule{})
	client := &domain.Client{ID: "c1", FullName: "Test", Email: "a@b.com", Country: "US", BirthDate: time.Now()}
	bank := &domain.Bank{ID: "b1", Name: "Gov Bank", Type: domain.BankTypeGovernment}
	result, err := engine.Evaluate(context.Background(), &EligibilityInput{
		Client:     client,
		Bank:       bank,
		MinPayment: 100,
		MaxPayment: 500,
		TermMonths: 24,
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Approved)
	assert.Equal(t, "PaymentRangeRule", result.RuleName)
}

func TestBankTypeRule_Government(t *testing.T) {
	r := BankTypeRule{}
	approved, priority, score := r.Evaluate(context.Background(), &EligibilityInput{
		Bank: &domain.Bank{Type: domain.BankTypeGovernment},
	})
	assert.True(t, approved)
	assert.Equal(t, 8, priority)
	assert.Equal(t, 0.8, score)
}

func TestBankTypeRule_Private(t *testing.T) {
	r := BankTypeRule{}
	approved, priority, score := r.Evaluate(context.Background(), &EligibilityInput{
		Bank: &domain.Bank{Type: domain.BankTypePrivate},
	})
	assert.True(t, approved)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 0.5, score)
}
