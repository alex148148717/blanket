package property_transactions_bl

import (
	"context"
	"property_transactions/property_transactions/property_transactions_db"
)

type DBClient interface {
	Add(ctx context.Context, userID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) error
}
type Client struct {
	dbClient DBClient
}

func New(dbClient DBClient) (*Client, error) {
	c := Client{dbClient: dbClient}
	return &c, nil
}

func (c *Client) Add(ctx context.Context, userID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error) {

	err := c.dbClient.Add(ctx, userID, txID, propertyTransactions)
	if err != nil {
		return 0, err
	}
	return txID, nil
}
