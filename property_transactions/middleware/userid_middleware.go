package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
)

const UserId = "user_id_key"

func UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if !IsValidUUID(userID) {
			//add more logic for userID
			response := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    1001,
					"message": "invalid userID",
					"user_id": userID,
				},
			}

			_ = json.NewEncoder(w).Encode(response)
			return
		}

		ctx := context.WithValue(r.Context(), UserId, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetUserID(ctx context.Context) string {
	userIDVal := ctx.Value(UserId)
	//no need to check  exist id
	return userIDVal.(string)
}
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
