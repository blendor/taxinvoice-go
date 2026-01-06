package model

import "time"

type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type TaxRate struct {
	ID            int64     `json:"id"`
	State         string    `json:"state"`
	Rate          float64   `json:"rate"`
	EffectiveDate time.Time `json:"effective_date"`
}

type Invoice struct {
	ID         int64         `json:"id"`
	CustomerID int64         `json:"customer_id"`
	State      string        `json:"state"`
	Subtotal   float64       `json:"subtotal"`
	TaxAmount  float64       `json:"tax_amount"`
	Total      float64       `json:"total"`
	Status     string        `json:"status"`
	DueDate    time.Time     `json:"due_date"`
	PaidAt     *time.Time    `json:"paid_at,omitempty"`
	Items      []InvoiceItem `json:"items,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
}

const (
	StatusUnpaid   = "unpaid"
	StatusPaid     = "paid"
	StatusOverdue  = "overdue"
	StatusRefunded = "refunded"
)

type InvoiceItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type TaxRequest struct {
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	State     string `json:"state"`
}

type TaxResponse struct {
	Subtotal  float64 `json:"subtotal"`
	TaxAmount float64 `json:"tax_amount"`
	Total     float64 `json:"total"`
}

type InvoiceRequest struct {
	CustomerID int64 `json:"customer_id"`
	State      string `json:"state"`
	Items      []struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	} `json:"items"`
}

type TaxReport struct {
	Period   Period          `json:"period"`
	ByState  []StateTax      `json:"by_state"`
	Totals   TaxReportTotals `json:"totals"`
}

type Period struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type StateTax struct {
	State        string  `json:"state"`
	TaxableSales float64 `json:"taxable_sales"`
	TaxCollected float64 `json:"tax_collected"`
	InvoiceCount int     `json:"invoice_count"`
}

type TaxReportTotals struct {
	TaxableSales float64 `json:"taxable_sales"`
	TaxCollected float64 `json:"tax_collected"`
	InvoiceCount int     `json:"invoice_count"`
}

type AgingReport struct {
	AsOf    time.Time    `json:"as_of"`
	Buckets []AgingBucket `json:"buckets"`
	Total   AgingTotal   `json:"total"`
}

type AgingBucket struct {
	Label        string    `json:"label"` // "current", "1-30", "31-60", "61-90", "90+"
	InvoiceCount int       `json:"invoice_count"`
	Amount       float64   `json:"amount"`
	Invoices     []Invoice `json:"invoices,omitempty"`
}

type AgingTotal struct {
	InvoiceCount int     `json:"invoice_count"`
	Amount       float64 `json:"amount"`
}
