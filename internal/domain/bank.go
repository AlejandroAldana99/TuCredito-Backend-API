package domain

// Type of banks (private or government).
type BankType string

const (
	BankTypePrivate    BankType = "PRIVATE"
	BankTypeGovernment BankType = "GOVERNMENT"
)

// Internal back structure
type Bank struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Type BankType `json:"type"`
}

// Input data required to create a bank
type CreateBankInput struct {
	Name string   `json:"name"`
	Type BankType `json:"type"`
}
