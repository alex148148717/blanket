package property_transactions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"property_transactions/property_transactions/property_transactions_bl"
	"property_transactions/property_transactions/property_transactions_db"
	"strconv"
	"sync"
	"testing"
	"time"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// test case
func TestInvalidUserID(t *testing.T) {
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

	addUrl := fmt.Sprintf("%s/property_transactions/v1/a/", server.URL)

	propertyTransactionsRequest := PropertyTransactionsRequest{PropertyID: "1", Amount: 100, Date: time.Now().Unix()}
	b, _ := json.Marshal(propertyTransactionsRequest)
	resp, err := http.Post(addUrl, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		t.Fatalf("json.Unmarshal: %s", string(body))
	}

	if errorResponse.Error.Code != 1001 {
		t.Errorf("Expected %d, got %d", 1001, errorResponse.Error.Code)
	}
}

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

	userID := uuid.NewString()
	propertyID := uuid.NewString()
	addUrl := fmt.Sprintf("%s/property_transactions/v1/user/%s/", server.URL, userID)

	addPropertyTransactions := func(t *testing.T, propertyTransactionsRequest PropertyTransactionsRequest) {
		b, _ := json.Marshal(propertyTransactionsRequest)
		resp, err := http.Post(addUrl, "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("Failed to send GET: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var propertyTransactionsResponse PropertyTransactionsResponse
		if err := json.Unmarshal(body, &propertyTransactionsResponse); err != nil {
			t.Fatalf("json.Unmarshal: %s", string(body))
		}
		if propertyTransactionsResponse.Success != true {
			t.Errorf("Expected %v, got %v", true, propertyTransactionsResponse.Success)
		}
	}
	addPropertyTransactions(t, PropertyTransactionsRequest{PropertyID: propertyID, Amount: 100, Date: time.Now().Unix()})
	addPropertyTransactions(t, PropertyTransactionsRequest{PropertyID: propertyID, Amount: -100, Date: time.Now().Unix()})
	addPropertyTransactions(t, PropertyTransactionsRequest{PropertyID: propertyID, Amount: 12.5, Date: time.Now().Unix()})

	allUrl := fmt.Sprintf("%s/property_transactions/v1/user/%s/property/%s/", server.URL, userID, propertyID)

	allPropertyTransactions := func(t *testing.T, queryParams QueryParams) {
		url := buildURL(allUrl, queryParams)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Failed to send GET: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var propertyTransactionsResponse PropertyTransactionsResponse
		if err := json.Unmarshal(body, &propertyTransactionsResponse); err != nil {
			t.Fatalf("json.Unmarshal: %s", string(body))
		}
		if propertyTransactionsResponse.Success != true {
			t.Errorf("Expected %v, got %v", true, propertyTransactionsResponse.Success)
		}
	}

	allPropertyTransactions(t, QueryParams{
		Type: "income", From: time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour).Unix(),
		To:   time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).Unix(),
		Page: 1, Limit: 20,
	})

	allPropertyTransactions(t, QueryParams{
		Type: "expense", From: time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour).Unix(),
		To:   time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour).Unix(),
		Page: 1, Limit: 20,
	})

}

type QueryParams struct {
	Type  string
	From  int64
	To    int64
	Page  int
	Limit int
}

func buildURL(baseURL string, params QueryParams) string {
	q := url.Values{}
	if params.Type != "" {
		q.Set("type", params.Type)
	}
	if params.From > 0 {
		q.Set("from", strconv.FormatInt(params.From, 10))
	}
	if params.To > 0 {
		q.Set("to", strconv.FormatInt(params.To, 10))
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	return fmt.Sprintf("%s?%s", baseURL, q.Encode())
}

func TestBalance(t *testing.T) {
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

	userID := uuid.NewString()
	propertyID := uuid.NewString()
	addUrl := fmt.Sprintf("%s/property_transactions/v1/user/%s/", server.URL, userID)

	addPropertyTransactions := func(t *testing.T, propertyTransactionsRequest PropertyTransactionsRequest) {
		b, _ := json.Marshal(propertyTransactionsRequest)
		resp, err := http.Post(addUrl, "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("Failed to send GET: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var propertyTransactionsResponse PropertyTransactionsResponse
		if err := json.Unmarshal(body, &propertyTransactionsResponse); err != nil {
			t.Fatalf("json.Unmarshal: %s", string(body))
		}
		if propertyTransactionsResponse.Success != true {
			t.Errorf("Expected %v, got %v", true, propertyTransactionsResponse.Success)
		}
	}

	var expectedBalance float64

	rand.Seed(time.Now().UnixNano())

	limit := 100
	workers := 20
	wg := sync.WaitGroup{}
	wg.Add(workers)
	balanceCHan := make(chan float64, limit)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for balance := range balanceCHan {
				addPropertyTransactions(t, PropertyTransactionsRequest{PropertyID: propertyID, Amount: balance, Date: time.Now().Unix()})
			}
		}()
	}
	for i := 0; i < limit; i++ {
		x := rand.Intn(2001) - 1000
		expectedBalance += float64(x)
		balanceCHan <- float64(x)
	}
	close(balanceCHan)
	wg.Wait()

	balanceURL := fmt.Sprintf("%s/property_transactions/v1/user/%s/property/%s/balance/", server.URL, userID, propertyID)

	resp, err := http.Get(balanceURL)
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var getPropertyBalanceHandlerResponse GetPropertyBalanceHandlerResponse
	if err := json.Unmarshal(body, &getPropertyBalanceHandlerResponse); err != nil {
		t.Fatalf("json.Unmarshal: %s", string(body))
	}
	if getPropertyBalanceHandlerResponse.Data.Balance != expectedBalance {
		t.Errorf("Expected %v, got %v", expectedBalance, getPropertyBalanceHandlerResponse.Data.Balance)
	}

}

func TestMonthlyReport(t *testing.T) {
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

	userID := uuid.NewString()
	propertyID := uuid.NewString()
	addUrl := fmt.Sprintf("%s/property_transactions/v1/user/%s/", server.URL, userID)

	addPropertyTransactions := func(t *testing.T, propertyTransactionsRequest PropertyTransactionsRequest) {
		b, _ := json.Marshal(propertyTransactionsRequest)
		resp, err := http.Post(addUrl, "application/json", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("Failed to send GET: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var propertyTransactionsResponse PropertyTransactionsResponse
		if err := json.Unmarshal(body, &propertyTransactionsResponse); err != nil {
			t.Fatalf("json.Unmarshal: %s", string(body))
		}
		if propertyTransactionsResponse.Success != true {
			t.Errorf("Expected %v, got %v", true, propertyTransactionsResponse.Success)
		}
	}

	var expectedBalance float64
	balance := []float64{100, 200, -50, 60, 30, -49, 4000, -34, 455, -945, 34, 56}

	for _, b := range balance {
		addPropertyTransactions(t, PropertyTransactionsRequest{PropertyID: propertyID, Amount: b, Date: time.Now().Unix()})
		expectedBalance += b
	}

	balanceURL := fmt.Sprintf("%s/property_transactions/v1/user/%s/property/%s/monthly_report/", server.URL, userID, propertyID)

	resp, err := http.Get(balanceURL)
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var getPropertyMonthlyReportResponse GetPropertyMonthlyReportResponse
	if err := json.Unmarshal(body, &getPropertyMonthlyReportResponse); err != nil {
		t.Fatalf("json.Unmarshal: %s", string(body))
	}
	balanceX := getPropertyMonthlyReportResponse.Data.MonthlyBalanceData.EndCash
	if balanceX != expectedBalance {
		t.Errorf("Expected %v, got %v", expectedBalance, balanceX)
	}

}
