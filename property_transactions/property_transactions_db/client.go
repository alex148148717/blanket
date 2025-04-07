package property_transactions_db

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	clickhouseConn clickhouse.Conn
}
type Config struct {
	ClickhouseOptions clickhouse.Options
}

func New(ctx context.Context, config Config) (*Client, error) {

	clickhouseConn, err := clickhouse.Open(&config.ClickhouseOptions)
	if err != nil {
		return nil, err
	}

	c := Client{clickhouseConn: clickhouseConn}
	return &c, nil
}

func (c *Client) Add(ctx context.Context, userID int, txID int, propertyTransactions PropertyTransactions) error {

	query := `
		INSERT INTO property_transactions 
		(id, user_id, property_id, amount, date, created_at)
		VALUES (?, ?, ?, ?, ?, now())
	`

	return c.clickhouseConn.Exec(ctx, query,
		txID,
		uint32(userID),
		propertyTransactions.PropertyID,
		propertyTransactions.Amount,
		propertyTransactions.Date,
	)

}
