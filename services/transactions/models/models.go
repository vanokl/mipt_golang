package models

type TransactionType string

const (
	Income   TransactionType = "income"
	Expense  TransactionType = "expense"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"userid"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Type        TransactionType `json:"type"`
	Category    string          `json:"category"`
	Date        string          `json:"date"`
	Description string          `json:"description"`
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	Ok          bool        `json:"ok"`
}

type ListResponse struct {
	Transaction []Transaction `json:"transaction"`
	Ok          bool          `json:"ok"`
}

type Commission struct {
	TransactionID   string          `json:"transaction_id"`
	Amount          float64         `json:"amount"`
	Currency        string          `json:"currency"`
	Type            TransactionType `json:"type"`
	Commission      float64         `json:"commission"`
	CalculationDate string          `json:"calculation_date"`
	Description     string          `json:"description"`
}
