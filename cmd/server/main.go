package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/blendor/taxinvoice-go/internal/handler"
	"github.com/blendor/taxinvoice-go/internal/service"
	"github.com/blendor/taxinvoice-go/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbURL := getenv("DATABASE_URL", "postgres://localhost/taxinvoice?sslmode=disable")
	port := getenv("PORT", "8080")

	db, err := store.New(dbURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := service.New(db)
	h := handler.New(svc, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /api/v1/calculate-tax", h.CalculateTax)
	mux.HandleFunc("POST /api/v1/invoices", h.CreateInvoice)
	mux.HandleFunc("GET /api/v1/invoices", h.ListInvoices)
	mux.HandleFunc("GET /api/v1/invoices/{id}", h.GetInvoice)
	mux.HandleFunc("PATCH /api/v1/invoices/{id}/pay", h.MarkPaid)
	mux.HandleFunc("PATCH /api/v1/invoices/{id}/refund", h.MarkRefunded)
	mux.HandleFunc("GET /api/v1/reports/tax", h.TaxReport)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      logging(log, mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("server shutting down")
	srv.Shutdown(ctx)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func logging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
