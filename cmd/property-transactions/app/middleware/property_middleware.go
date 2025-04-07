package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

const PropertyId = "property_id_key"

func PropertyIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		propertyId := chi.URLParam(r, "propertyID")
		ctx := context.WithValue(r.Context(), PropertyId, propertyId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetPropertyId(ctx context.Context) (string, error) {
	userIDVal := ctx.Value(PropertyId)
	if userIDVal == nil {
		return "", fmt.Errorf("property ID not found in context")
	}
	propertyID, ok := userIDVal.(string)
	if !ok {
		return "", fmt.Errorf("property ID is not of type int")
	}

	return propertyID, nil
}
