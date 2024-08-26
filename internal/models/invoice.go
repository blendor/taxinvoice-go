package models

import (
	"context"
	"time"
)

type Invoice struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	TaxAmount   float64   `json:"tax_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []InvoiceItem `json:"items"`
}

type InvoiceItem struct {
	ID        int64   `json:"id"`
	InvoiceID int64   `json:"invoice_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type InvoiceGenerationRequest struct {
	CustomerID int64         `json:"customer_id"`
	Items      []InvoiceItem `json:"items"`
	State      string        `json:"state"`
}

type InvoiceGenerationResponse struct {
	Invoice     Invoice `json:"invoice"`
	PDFLocation string  `json:"pdf_location"`
	CSVLocation string  `json:"csv_location"`
}

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoice(ctx context.Context, id int64) (*Invoice, error)
	ListInvoices(ctx context.Context, customerID int64, limit, offset int) ([]*Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	DeleteInvoice(ctx context.Context, id int64) error
}