package api

import (
	"github.com/gorilla/mux"
	"github.com/blendor/taxinvoice-go/internal/api/handlers"
	"github.com/blendor/taxinvoice-go/internal/api/middleware"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

func SetupRoutes(r *mux.Router, taxHandler *handlers.TaxHandler, invoiceHandler *handlers.InvoiceHandler, logger *logger.Logger) {
	// Apply global middleware
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Auth(logger))

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Tax routes
	api.HandleFunc("/calculate-tax", taxHandler.CalculateTax).Methods("POST")

	// Invoice routes
	api.HandleFunc("/generate-invoice", invoiceHandler.GenerateInvoice).Methods("POST")

	// Health check route
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
}