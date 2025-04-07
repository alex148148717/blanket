package property_transactions

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"property_transactions/property_transactions/property_transactions_db"
	"time"
)

type PropertyTransactionsClient interface {
	Add(ctx context.Context, userID int, txID int, propertyTransactions property_transactions_db.PropertyTransactions) (int, error)
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
					middleware.Compress(5),
					middleware.Recoverer,
				)
				r.Post("/", s.addPropertyTransactionsHandler)
				//r.Get("/all", s.getPropertyTransactionsHandler)
				//r.Get("/property/{propertyID}/", s.getPropertyBalanceHandler)
				//r.Get("/property/{propertyID}/monthly_report", s.addPropertyTransactionsHandler)

			})
		})
	})

	return &s, nil
}

func (s *Server) addPropertyTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx
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

	_, err = s.propertyTransactionsClient.Add(ctx, 1, txID, propertyTransactions)

	w.WriteHeader(http.StatusCreated)

	return
}

func (s *Server) getPropertyTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx
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
