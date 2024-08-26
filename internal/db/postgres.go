package db

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func NewPostgresConnection(connString string) (*sql.DB, error) {
    return sql.Open("postgres", connString)
}