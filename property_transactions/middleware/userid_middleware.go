package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const UserId = "user_id_key"

func UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			//add more logic for userID
			response := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    1001,
					"message": "invalid userID",
					"user_id": userIDStr,
				},
			}

			_ = json.NewEncoder(w).Encode(response)
			return
		}

		ctx := context.WithValue(r.Context(), UserId, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetUserID(ctx context.Context) int {
	userIDVal := ctx.Value(UserId)
	//no need to check  exist id
	return userIDVal.(int)
}
