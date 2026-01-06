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

	now := time.Now()
	inv.CreatedAt = now
	inv.Status = model.StatusUnpaid
	inv.DueDate = now.AddDate(0, 0, 30) // Net 30

	err = tx.QueryRowContext(ctx,
		`INSERT INTO invoices (customer_id, state, subtotal, tax_amount, total, status, due_date, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		inv.CustomerID, inv.State, inv.Subtotal, inv.TaxAmount, inv.Total, inv.Status, inv.DueDate, inv.CreatedAt,
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

func (s *Store) GetInvoice(ctx context.Context, id int64) (*model.Invoice, error) {
	var inv model.Invoice
	var paidAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, customer_id, state, subtotal, tax_amount, total, status, due_date, paid_at, created_at
		 FROM invoices WHERE id = $1`, id,
	).Scan(&inv.ID, &inv.CustomerID, &inv.State, &inv.Subtotal, &inv.TaxAmount, &inv.Total,
		&inv.Status, &inv.DueDate, &paidAt, &inv.CreatedAt)
	if err != nil {
		return nil, err
	}
	if paidAt.Valid {
		inv.PaidAt = &paidAt.Time
	}
	return &inv, nil
}

func (s *Store) ListInvoices(ctx context.Context, status string) ([]model.Invoice, error) {
	query := `SELECT id, customer_id, state, subtotal, tax_amount, total, status, due_date, paid_at, created_at
			  FROM invoices`
	args := []any{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var inv model.Invoice
		var paidAt sql.NullTime
		err := rows.Scan(&inv.ID, &inv.CustomerID, &inv.State, &inv.Subtotal, &inv.TaxAmount, &inv.Total,
			&inv.Status, &inv.DueDate, &paidAt, &inv.CreatedAt)
		if err != nil {
			return nil, err
		}
		if paidAt.Valid {
			inv.PaidAt = &paidAt.Time
		}
		invoices = append(invoices, inv)
	}
	return invoices, rows.Err()
}

func (s *Store) UpdateInvoiceStatus(ctx context.Context, id int64, status string, paidAt *time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE invoices SET status = $1, paid_at = $2 WHERE id = $3`,
		status, paidAt, id,
	)
	return err
}

func (s *Store) GetTaxReport(ctx context.Context, from, to time.Time) ([]model.StateTax, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT state, SUM(subtotal) as taxable_sales, SUM(tax_amount) as tax_collected, COUNT(*) as invoice_count
		 FROM invoices
		 WHERE created_at >= $1 AND created_at < $2 AND status != $3
		 GROUP BY state
		 ORDER BY tax_collected DESC`,
		from, to, model.StatusRefunded,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.StateTax
	for rows.Next() {
		var st model.StateTax
		if err := rows.Scan(&st.State, &st.TaxableSales, &st.TaxCollected, &st.InvoiceCount); err != nil {
			return nil, err
		}
		results = append(results, st)
	}
	return results, rows.Err()
}
