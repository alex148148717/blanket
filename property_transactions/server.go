package property_transactions

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	middleware2 "property_transactions/property_transactions/middleware"
	"property_transactions/property_transactions/property_transactions_bl"
	"property_transactions/property_transactions/property_transactions_db"
	"strconv"
	"time"
)

type PropertyTransactionsClient interface {
	Add(ctx context.Context, userID int, propertyID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error)
	All(ctx context.Context, userID int, propertyId int, propertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error)
	Balance(ctx context.Context, userID int, propertyID int) (float64, error)
	MonthlyBalance(ctx context.Context, userID int, propertyID int, from time.Time, to time.Time) (*property_transactions_bl.MonthlyBalanceData, error)
}

type Server struct {
	propertyTransactionsClient PropertyTransactionsClient
}

func New(ctx context.Context, router *chi.Mux, propertyTransactionsClient PropertyTransactionsClient) (*Server, error) {
	s := Server{propertyTransactionsClient}

	router.Route("/property_transactions/v1/", func(r chi.Router) {
		r.Route("/user/{userID}/", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(
					middleware2.UserIDMiddleware,
					middleware2.PropertyIDMiddleware,
					middleware.Recoverer,
				)
				r.Post("/", s.addPropertyTransactionsHandler)
				r.Get("/property/{propertyID}/", s.getPropertyTransactionsHandler)
				r.Get("/property/{propertyID}/balance/", s.getPropertyBalanceHandler)
				r.Get("/property/{propertyID}/monthly_report/", s.getPropertyMonthlyReportHandler)

			})
		})
	})

	return &s, nil
}

func (s *Server) addPropertyTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware2.GetUserID(ctx)

	var req PropertyTransactionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	//todo ...
	txID := int(time.Now().UnixNano())
	propertyTransactions, err := req.ToModel()
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = s.propertyTransactionsClient.Add(ctx, userID, propertyTransactions.PropertyID, txID, propertyTransactions)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}
	_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: true})

	w.WriteHeader(http.StatusCreated)

	return
}

func (s *Server) getPropertyTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware2.GetUserID(ctx)
	propertyId, err := middleware2.GetPropertyId(ctx)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}
	getPropertyTransactionsRequest, err := parseGetPropertyTransactionsRequest(r)

	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}

	transactions, err := s.propertyTransactionsClient.All(ctx, userID, propertyId, *getPropertyTransactionsRequest)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}
	resTransactions := make([]Transaction, 0, len(transactions))
	for _, t := range transactions {
		resTransactions = append(resTransactions, Transaction{t.Amount, t.Date.Unix()})
	}
	res := GetPropertyTransactionsHandlerResponse{
		Success: true,
		Data:    TransactionList{Transactions: resTransactions},
	}
	_ = json.NewEncoder(w).Encode(res)

	w.WriteHeader(http.StatusCreated)
	return

}

func (s *Server) getPropertyBalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware2.GetUserID(ctx)
	propertyId, err := middleware2.GetPropertyId(ctx)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}
	balance, err := s.propertyTransactionsClient.Balance(ctx, userID, propertyId)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}

	_ = json.NewEncoder(w).Encode(GetPropertyBalanceHandlerResponse{Success: true, Data: Balance{balance}})

	w.WriteHeader(http.StatusCreated)

	return

}

func (s *Server) getPropertyMonthlyReportHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware2.GetUserID(ctx)
	propertyId, err := middleware2.GetPropertyId(ctx)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}

	q := r.URL.Query()

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)

	if fromUnix, err := strconv.ParseInt(q.Get("from"), 10, 64); err == nil && fromUnix > 0 {
		from = time.Unix(fromUnix, 0)
	}

	if toUnix, err := strconv.ParseInt(q.Get("to"), 10, 64); err == nil && toUnix > 0 {
		to = time.Unix(toUnix, 0)
	}

	balance, err := s.propertyTransactionsClient.MonthlyBalance(ctx, userID, propertyId, from, to)
	if err != nil {
		_ = json.NewEncoder(w).Encode(PropertyTransactionsResponse{Success: false, Error: ErrorDetail{Code: 1002}})
		return
	}
	_ = balance
	_ = json.NewEncoder(w).Encode(GetPropertyMonthlyReportResponse{Success: true, Data: MonthlyReport{ConvertMonthlyBalanceData(*balance)}})

	w.WriteHeader(http.StatusCreated)

	return

}
