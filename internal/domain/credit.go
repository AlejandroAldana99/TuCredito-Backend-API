package domain

import "time"

// Status of a credit application (PENDING, APPROVED, REJECTED)
type CreditStatus string

const (
	CreditStatusPending  CreditStatus = "PENDING"
	CreditStatusApproved CreditStatus = "APPROVED"
	CreditStatusRejected CreditStatus = "REJECTED"
)

// Type of credit product (AUTO, MORTGAGE, COMMERCIAL)
type CreditType string

const (
	CreditTypeAuto       CreditType = "AUTO"
	CreditTypeMortgage   CreditType = "MORTGAGE"
	CreditTypeCommercial CreditType = "COMMERCIAL"
)

// Credit structure
type Credit struct {
	ID         string       `json:"id"`
	ClientID   string       `json:"client_id"`
	BankID     string       `json:"bank_id"`
	MinPayment float64      `json:"min_payment"`
	MaxPayment float64      `json:"max_payment"`
	TermMonths int          `json:"term_months"`
	CreditType CreditType   `json:"credit_type"`
	CreatedAt  time.Time    `json:"created_at"`
	Status     CreditStatus `json:"status"`
	IsActive   bool         `json:"is_active"`
}

// Structure for creating a credit
type CreateCreditInput struct {
	ClientID   string     `json:"client_id"`
	BankID     string     `json:"bank_id"`
	MinPayment float64    `json:"min_payment"`
	MaxPayment float64    `json:"max_payment"`
	TermMonths int        `json:"term_months"`
	CreditType CreditType `json:"credit_type"`
}

// Structure for updating a credit status
type UpdateCreditStatusInput struct {
	Status CreditStatus `json:"status"`
}

// Structure for updating a credit
type UpdateCreditInput struct {
	MinPayment float64      `json:"min_payment"`
	MaxPayment float64      `json:"max_payment"`
	TermMonths int          `json:"term_months"`
	Status     CreditStatus `json:"status"`
}
