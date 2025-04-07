package property_transactions_db

import "time"

type PropertyTransactions struct {
	PropertyID int
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
	UserID     uint32    `json:"user_id"`
	PropertyID uint32    `json:"property_id"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}
