package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const PropertyId = "property_id_key"

func PropertyIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		propertyId, _ := strconv.Atoi(chi.URLParam(r, "propertyID"))
		ctx := context.WithValue(r.Context(), PropertyId, propertyId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetPropertyId(ctx context.Context) (int, error) {
	userIDVal := ctx.Value(PropertyId)
	if userIDVal == nil {
		return 0, fmt.Errorf("property ID not found in context")
	}
	propertyID, ok := userIDVal.(int)
	if !ok {
		return 0, fmt.Errorf("property ID is not of type int")
	}

	return propertyID, nil
}
