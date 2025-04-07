package property_transactions

import (
	"errors"
	"property_transactions/property_transactions/property_transactions_db"
	"time"
)

type PropertyTransactionsRequest struct {
	PropertyID int     `json:"propertyID"`
	Amount     float64 `json:"amount"`
	Date       int64   `json:"date"`
}

func (r *PropertyTransactionsRequest) ToModel() (property_transactions_db.PropertyTransactions, error) {

	ret := property_transactions_db.PropertyTransactions{
		PropertyID: r.PropertyID,
		Amount:     r.Amount,
		Date:       time.Unix(int64(r.Date), 0),
	}
	if r.Date == 0 {
		return ret, errors.New("date is required")
	}
	return ret, nil
}

type PropertyTransactionsResponse struct {
}
