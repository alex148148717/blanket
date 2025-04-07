package property_transactions_bl

import (
	"context"
	"property_transactions/property_transactions/property_transactions_db"
	"time"
)

type DBClient interface {
	Add(ctx context.Context, userID string, propertyID string, txID int, propertyTransactions property_transactions_db.PropertyTransactions) error
	All(ctx context.Context, userID string, propertyID string, allPropertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error)
	Balance(ctx context.Context, userID string, propertyID string) (float64, error)
	MonthlyBalance(ctx context.Context, userID string, propertyID string, from time.Time, to time.Time) ([]property_transactions_db.Transaction, error)
}
type Client struct {
	dbClient DBClient
}

func New(dbClient DBClient) (*Client, error) {
	c := Client{dbClient: dbClient}
	return &c, nil
}

func (c *Client) Add(ctx context.Context, userID string, propertyID string, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error) {

	err := c.dbClient.Add(ctx, userID, propertyID, txID, propertyTransactions)
	if err != nil {
		return 0, err
	}
	return txID, nil
}
func (c *Client) All(ctx context.Context, userID string, propertyID string, propertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error) {

	transactions, err := c.dbClient.All(ctx, userID, propertyID, propertyTransactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func (c *Client) Balance(ctx context.Context, userID string, propertyID string) (float64, error) {
	return c.dbClient.Balance(ctx, userID, propertyID)

}

func (c *Client) MonthlyBalance(ctx context.Context, userID string, propertyID string, from time.Time, to time.Time) (*MonthlyBalanceData, error) {
	transactions, err := c.dbClient.MonthlyBalance(ctx, userID, propertyID, from, to)
	if err != nil {
		return nil, err
	}

	_ = transactions
	monthlyBalanceData := MonthlyBalanceData{}

	if len(transactions) == 0 {
		return &monthlyBalanceData, nil
	}
	monthlyBalanceData.StartingCash = transactions[0].Amount
	monthlyBalanceData.EndCash += transactions[0].Amount
	records := make([]Record, 0, len(transactions)-1)

	for i := 1; i < len(transactions); i++ {
		transaction := transactions[i]
		monthlyBalanceData.EndCash += transaction.Amount
		transactionType := property_transactions_db.TransactionTypeIncome
		if transaction.Amount < 0 {
			transactionType = property_transactions_db.TransactionTypeExpense
		}
		records = append(records, Record{
			Record:          i,
			TransactionType: transactionType,
			Amount:          transaction.Amount,
			Total:           monthlyBalanceData.EndCash,
		})
	}
	monthlyBalanceData.Records = records
	return &monthlyBalanceData, nil
}
