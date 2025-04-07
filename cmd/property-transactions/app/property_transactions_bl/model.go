package property_transactions_bl

import "property_transactions/cmd/property-transactions/app/property_transactions_db"

type Record struct {
	Record          int
	TransactionType property_transactions_db.TransactionType
	Amount          float64
	Total           float64
}

type MonthlyBalanceData struct {
	StartingCash float64
	Records      []Record
	EndCash      float64
}
