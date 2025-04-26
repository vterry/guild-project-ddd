package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/app/token"
	"github.com/vterry/ddd-study/auth-server/internal/app/utils"
)

func Auhtentication() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				denied(w)
				return
			}

			fields := strings.Fields(authHeader)
			if len(fields) != 2 || fields[0] != "Bearer" {
				denied(w)
				return
			}

			jwtToken := fields[1]

			claims, err := token.ValidateJWT(jwtToken)
			if err != nil {
				denied(w)
				return
			}

			if time.Now().After(claims.ExpiresAt.Time) {
				denied(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func denied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
