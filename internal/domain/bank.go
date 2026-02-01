package domain

// Type of banks (private or government).
type BankType string

const (
	BankTypePrivate    BankType = "PRIVATE"
	BankTypeGovernment BankType = "GOVERNMENT"
)

// Bank structure
type Bank struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Type     BankType `json:"type"`
	IsActive bool     `json:"is_active"`
}

// Structure for creating a bank
type CreateBankInput struct {
	Name string   `json:"name"`
	Type BankType `json:"type"`
}

// Structure for updating a bank
type UpdateBankInput struct {
	Name string   `json:"name"`
	Type BankType `json:"type"`
}
