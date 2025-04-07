package property_transactions_db

import "time"

type PropertyTransactions struct {
	PropertyID string
	Amount     float64
	Date       time.Time
}
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type AllPropertyTransactionsParams struct {
	Type  TransactionType
	From  time.Time
	TO    time.Time
	Page  int
	Limit int
}

type Transaction struct {
	UserID     string    `json:"user_id"`
	PropertyID string    `json:"property_id"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}
