package middleware

import (
	"net/http"
	"strings"

	"github.com/vterry/ddd-study/character/internal/adapters/input/token"
)

func Auhtentication(tokenAdapter token.TokenValidationAdapter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header has the correct format
			parts := strings.Fields(authHeader)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			valid, err := tokenAdapter.TokenValidation(r.Context(), token)
			if err != nil || !valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
