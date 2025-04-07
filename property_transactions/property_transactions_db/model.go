package property_transactions_db

import "time"

type PropertyTransactions struct {
	PropertyID int
	Amount     float64
	Date       time.Time
}
