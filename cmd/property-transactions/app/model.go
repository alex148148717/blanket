package app

import (
	"errors"
	"net/http"
	"property_transactions/cmd/property-transactions/app/property_transactions_bl"
	"property_transactions/cmd/property-transactions/app/property_transactions_db"

	"strconv"
	"time"
)

type PropertyTransactionsRequest struct {
	PropertyID string  `json:"propertyID"`
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
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func parseTransactionType(s string) property_transactions_db.TransactionType {
	switch property_transactions_db.TransactionType(s) {
	case property_transactions_db.TransactionTypeIncome:
		return property_transactions_db.TransactionTypeIncome
	default:
		return property_transactions_db.TransactionTypeIncome
	}
}

func parseGetPropertyTransactionsRequest(r *http.Request) (*property_transactions_db.AllPropertyTransactionsParams, error) {
	q := r.URL.Query()

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)

	fromUnix, err := strconv.ParseInt(q.Get("from"), 10, 64)
	if err == nil && fromUnix > 0 {
		from = time.Unix(fromUnix, 0)
	}

	toUnix, err := strconv.ParseInt(q.Get("to"), 10, 64)
	if err == nil && toUnix > 0 {
		to = time.Unix(toUnix, 0)
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	getPropertyTransactionsRequest := property_transactions_db.AllPropertyTransactionsParams{
		Type:  parseTransactionType(q.Get("type")),
		From:  from,
		TO:    to,
		Page:  page,
		Limit: limit,
	}
	return &getPropertyTransactionsRequest, nil
}

type Transaction struct {
	Amount float64 `json:"amount"`
	Date   int64   `json:"date"`
}

type TransactionList struct {
	PropertyID   string        `json:"propertyId"`
	UserID       string        `json:"userId"`
	Transactions []Transaction `json:"transactions"`
}

type GetPropertyTransactionsHandlerResponse struct {
	Success bool `json:"success"`

	Data TransactionList `json:"data"`
}

type Balance struct {
	PropertyID string `json:"propertyId"`
	UserID     string `json:"userId"`

	Balance float64 `json:"balance"`
}
type GetPropertyBalanceHandlerResponse struct {
	Success bool    `json:"success"`
	Data    Balance `json:"data"`
}

type Record struct {
	Record          int                                      `json:"id"`
	TransactionType property_transactions_db.TransactionType `json:"type"`
	Amount          float64                                  `json:"amount"`
	Total           float64                                  `json:"total"`
}

type MonthlyBalanceData struct {
	StartingCash float64  `json:"startingCash"`
	Records      []Record `json:"records"`
	EndCash      float64  `json:"endCash"`
}
type MonthlyReport struct {
	PropertyID string `json:"propertyID"`
	UserID     string `json:"userId"`

	MonthlyBalanceData MonthlyBalanceData `json:"monthlyBalance"`
}
type GetPropertyMonthlyReportResponse struct {
	Success bool          `json:"success"`
	Data    MonthlyReport `json:"data"`
}

func ConvertMonthlyBalanceData(monthlyBalanceData property_transactions_bl.MonthlyBalanceData) MonthlyBalanceData {
	records := make([]Record, 0, len(monthlyBalanceData.Records))
	for _, r := range monthlyBalanceData.Records {
		records = append(records, Record{
			Record:          r.Record,
			TransactionType: r.TransactionType,
			Amount:          r.Amount,
			Total:           r.Total,
		})
	}

	return MonthlyBalanceData{
		StartingCash: monthlyBalanceData.StartingCash,
		Records:      records,
		EndCash:      monthlyBalanceData.EndCash,
	}
}
