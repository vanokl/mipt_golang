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

type TransactionConverted struct {
	ID                string          `json:"id"`
	UserID            string          `json:"userid"`
	Amount            float64         `json:"amount"`
	Currency          string          `json:"currency"`
	ConvertedAmount   float64         `json:"converted_amount"`
	ConvertedCurrency string          `json:"converted_currency"`
	Type              TransactionType `json:"type"`
	Category          string          `json:"category"`
	Date              string          `json:"date"`
	Description       string          `json:"description"`
}

type DataCurrency struct {
	EUR float64 `json:"EUR"`
	RUB float64 `json:"RUB"`
	USD float64 `json:"USD"`
}

type CurrencyStruct struct {
	DataCurrency DataCurrency `json:"data"`
}

type CurrencyAnswer struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	OriginalAmount  float64 `json:"original_amount"`
	ConvertedAmount float64 `json:"converted_amount"`
	Rate            float64 `json:"rate"`
}
