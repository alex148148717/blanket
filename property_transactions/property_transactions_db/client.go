package property_transactions_db

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"strings"
	"time"
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

func (c *Client) Add(ctx context.Context, userID int, propertyID int, txID int, propertyTransactions PropertyTransactions) error {

	query := `
		INSERT INTO property_transactions 
		(id, user_id, property_id, amount, date, created_at)
		VALUES (?, ?, ?, ?, ?, now())
	`

	return c.clickhouseConn.Exec(ctx, query,
		txID,
		uint32(userID),
		propertyID,
		propertyTransactions.Amount,
		propertyTransactions.Date,
	)

}

func (c *Client) All(ctx context.Context, userID int, propertyID int, allPropertyTransactions AllPropertyTransactionsParams) ([]Transaction, error) {

	query := `
		SELECT  user_id, property_id, amount,  date
		FROM property_transactions
		WHERE user_id = ? AND property_id = ?
		  AND date >= toDate(?) AND date <= toDate(?)
	`

	args := []interface{}{
		userID,
		propertyID,
		allPropertyTransactions.From.Unix(),
		allPropertyTransactions.TO.Unix(),
	}

	switch allPropertyTransactions.Type {
	case TransactionTypeIncome:
		query += " AND amount > 0"
	case TransactionTypeExpense:
		query += " AND amount < 0"
	}

	query += `
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`

	offset := (allPropertyTransactions.Page - 1) * allPropertyTransactions.Limit
	args = append(args, allPropertyTransactions.Limit, offset)

	rows, err := c.clickhouseConn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Transaction
	for rows.Next() {
		var t Transaction

		if err := rows.Scan(&t.UserID, &t.PropertyID, &t.Amount, &t.Date); err != nil {
			return nil, err
		}

		results = append(results, t)
	}
	return results, nil

}
func FormatQuery(query string, args []interface{}) string {
	for _, arg := range args {
		var val string
		switch v := arg.(type) {
		case string:
			val = fmt.Sprintf("'%s'", v)
		case time.Time:
			val = fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
		default:
			val = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", val, 1)
	}
	return query
}

func (c *Client) Balance(ctx context.Context, userID int, propertyID int) (float64, error) {

	query := `
		SELECT  sum(amount)as amount
		FROM property_transactions
		WHERE user_id = ? AND property_id = ?
	`

	args := []interface{}{userID, propertyID}

	rows, err := c.clickhouseConn.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var balance float64
	if rows.Next() {
		_ = rows.Scan(&balance)
	}
	return balance, nil

}
