package property_transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"net/http/httptest"
	"property_transactions/property_transactions/property_transactions_bl"
	"property_transactions/property_transactions/property_transactions_db"
	"testing"
	"time"
)

func TestAddPropertyTransactions(t *testing.T) {
	ctx := context.Background()

	clickhouseOptions := clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "myuser",
			Password: "mypassword",
		},
		DialTimeout: time.Second * 10,
		Debug:       false,
	}

	propertyTransactionsDBClient, err := property_transactions_db.New(ctx, property_transactions_db.Config{ClickhouseOptions: clickhouseOptions})
	if err != nil {
		t.Error(err)
	}
	propertyTransactionsClient, err := property_transactions_bl.New(propertyTransactionsDBClient)
	if err != nil {
		t.Error(err)
	}
	r := chi.NewRouter()
	s, err := New(ctx, r, propertyTransactionsClient)
	if err != nil {
		t.Error(err)
	}
	_ = s
	server := httptest.NewServer(r)
	defer server.Close()

	addUrl := fmt.Sprintf("%s/property_transactions/v1/12/", server.URL)
	t.Log(addUrl)
	propertyTransactionsRequest := PropertyTransactionsRequest{
		PropertyID: 1,
		Amount:     100,
		Date:       time.Now().Unix(),
	}
	b, _ := json.Marshal(propertyTransactionsRequest)
	resp, err := http.Post(addUrl, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	want := "Hello from server"

	if got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}
