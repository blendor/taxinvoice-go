package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/blendor/taxinvoice-go/internal/api/handlers"
	"github.com/blendor/taxinvoice-go/internal/config"
	"github.com/blendor/taxinvoice-go/internal/db"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database connection
	database, err := db.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close()

	// Initialize router
	router := mux.NewRouter()

	// Initialize handlers
	taxHandler := handlers.NewTaxHandler(database, logger)
	invoiceHandler := handlers.NewInvoiceHandler(database, logger)

	// In main.go, replace the routing setup with:
	router := mux.NewRouter()
	api.SetupRoutes(router, taxHandler, invoiceHandler, logger)

	// Set up routes
	// router.HandleFunc("/api/v1/calculate-tax", taxHandler.CalculateTax).Methods("POST")
	// router.HandleFunc("/api/v1/generate-invoice", invoiceHandler.GenerateInvoice).Methods("POST")

	// Set up middleware
	router.Use(loggingMiddleware(logger))

	// Create server
	// srv := &http.Server{
	//		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
	//		WriteTimeout: time.Second * 15,
	//		ReadTimeout:  time.Second * 15,
	//		IdleTimeout:  time.Second * 60,
	//		Handler:      router,
	//	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Use the router with the set up routes
	}


// Create server
srv := &http.Server{
    Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
    WriteTimeout: time.Second * 15,
    ReadTimeout:  time.Second * 15,
    IdleTimeout:  time.Second * 60,
    Handler:      router, // Use the router with the set up routes
}

	// Start server
	go func() {
		logger.Info("Starting server", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Server error", "error", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	logger.Info("Server shutting down")
	os.Exit(0)
}

func loggingMiddleware(logger *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}