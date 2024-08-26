package middleware

import (
	"net/http"
	"strings"

	"github.com/blendor/taxinvoice-go/pkg/logger"
)

func Auth(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				logger.Warn("Missing authorization token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Remove 'Bearer ' prefix if present
			token = strings.TrimPrefix(token, "Bearer ")

			// TODO: Implement proper token validation
			if token != "your-secret-token" {
				logger.Warn("Invalid authorization token", "token", token)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}