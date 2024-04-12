package models

import "time"

type TransactionType string

const (
	Income   TransactionType = "income"
	Expense  TransactionType = "expense"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID          string          `json:"id"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Type        TransactionType `json:"type"`
	Category    string          `json:"category"`
	Date        time.Time       `json:"date"`
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
