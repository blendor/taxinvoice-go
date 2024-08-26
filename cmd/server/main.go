package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/yourusername/tax-invoice-app/internal/config"
    "github.com/yourusername/tax-invoice-app/internal/db"
)

func main() {
    cfg := config.LoadConfig()

    db, err := db.NewPostgresConnection(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    r := mux.NewRouter()

    // Add routes here

    // Add this import
    "github.com/yourusername/tax-invoice-app/internal/api/handlers"

    // In the main function, before starting the server:
    r.HandleFunc("/calculate-tax", handlers.CalculateTax).Methods("POST")

    log.Printf("Starting server on :%s", cfg.ServerPort)
    log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}