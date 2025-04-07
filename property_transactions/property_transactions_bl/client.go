package property_transactions_bl

import (
	"context"
	"property_transactions/property_transactions/property_transactions_db"
)

type DBClient interface {
	Add(ctx context.Context, userID int, propertyID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) error
	All(ctx context.Context, userID int, propertyID int, allPropertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error)
}
type Client struct {
	dbClient DBClient
}

func New(dbClient DBClient) (*Client, error) {
	c := Client{dbClient: dbClient}
	return &c, nil
}

func (c *Client) Add(ctx context.Context, userID int, propertyID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error) {

	err := c.dbClient.Add(ctx, userID, propertyID, txID, propertyTransactions)
	if err != nil {
		return 0, err
	}
	return txID, nil
}
func (c *Client) All(ctx context.Context, userID int, propertyID int, propertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error) {

	transactions, err := c.dbClient.All(ctx, userID, propertyID, propertyTransactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
