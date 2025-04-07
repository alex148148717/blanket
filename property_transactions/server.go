package property_transactions

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	middleware2 "property_transactions/property_transactions/middleware"
	"property_transactions/property_transactions/property_transactions_db"
	"time"
)

type PropertyTransactionsClient interface {
	Add(ctx context.Context, userID int, propertyID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error)
	All(ctx context.Context, userID int, propertyId int, propertyTransactions property_transactions_db.AllPropertyTransactionsParams) ([]property_transactions_db.Transaction, error)
}

type Server struct {
	propertyTransactionsClient PropertyTransactionsClient
}

func New(ctx context.Context, router *chi.Mux, propertyTransactionsClient PropertyTransactionsClient) (*Server, error) {
	s := Server{propertyTransactionsClient}

	router.Route("/property_transactions/v1/", func(r chi.Router) {
		r.Route("/{userID}/", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(
					middleware2.UserIDMiddleware,
					middleware2.PropertyIDMiddleware,
					middleware.Compress(5),
					middleware.Recoverer,
				)
				r.Post("/", s.addPropertyTransactionsHandler)
				r.Get("/{propertyID}/", s.getPropertyTransactionsHandler)
				//r.Get("/property/{propertyID}/", s.getPropertyBalanceHandler)
				//r.Get("/property/{propertyID}/monthly_report", s.addPropertyTransactionsHandler)

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
	_ = ctx
	return

}
func (s *Server) getPropertyMonthlyReportHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx
	return

}
