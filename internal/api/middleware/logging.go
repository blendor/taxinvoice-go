package middleware

import (
	"net/http"
	"time"

	"github.com/blendor/taxinvoice-go/pkg/logger"
)

func Logging(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture the status code
			wrappedWriter := wrapResponseWriter(w)

			// Call the next handler
			next.ServeHTTP(wrappedWriter, r)

			// Log the request details
			logger.Info("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrappedWriter.status,
				"duration", time.Since(start),
				"ip", r.RemoteAddr,
			)
		})
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK}
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}