package domain

import "time"

// lifecycle state of a credit application (PENDING, APPROVED, REJECTED)
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

// Structure of a credit application
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
}

// Data required to create a credit application
type CreateCreditInput struct {
	ClientID   string     `json:"client_id"`
	BankID     string     `json:"bank_id"`
	MinPayment float64    `json:"min_payment"`
	MaxPayment float64    `json:"max_payment"`
	TermMonths int        `json:"term_months"`
	CreditType CreditType `json:"credit_type"`
}

// Data to update credit status
type UpdateCreditStatusInput struct {
	Status CreditStatus `json:"status"`
}
