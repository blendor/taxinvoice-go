package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

type Database struct {
	*sql.DB
	Logger *logger.Logger
}

func NewPostgresConnection(connectionString string, logger *logger.Logger) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %v", err)
	}

	logger.Info("Successfully connected to database")

	return &Database{DB: db, Logger: logger}, nil
}

func (db *Database) Close() {
	if err := db.DB.Close(); err != nil {
		db.Logger.Error("Error closing database connection", "error", err)
	}
}