package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/blendor/taxinvoice-go/internal/model"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New(databaseURL string) (*Store, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	var p model.Product
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, price, category FROM products WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Category)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) GetTaxRate(ctx context.Context, state string) (*model.TaxRate, error) {
	var t model.TaxRate
	err := s.db.QueryRowContext(ctx,
		`SELECT id, state, rate, effective_date FROM tax_rates
		 WHERE state = $1 AND effective_date <= $2
		 ORDER BY effective_date DESC LIMIT 1`,
		state, time.Now(),
	).Scan(&t.ID, &t.State, &t.Rate, &t.EffectiveDate)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) CreateInvoice(ctx context.Context, inv *model.Invoice) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO invoices (customer_id, state, subtotal, tax_amount, total, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		inv.CustomerID, inv.State, inv.Subtotal, inv.TaxAmount, inv.Total, time.Now(),
	).Scan(&inv.ID)
	if err != nil {
		return err
	}

	for _, item := range inv.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO invoice_items (invoice_id, product_id, quantity, unit_price, subtotal)
			 VALUES ($1, $2, $3, $4, $5)`,
			inv.ID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
