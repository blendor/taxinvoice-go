package models

import (
	"context"
	"time"
)

type TaxRate struct {
	ID            int64     `json:"id"`
	State         string    `json:"state"`
	Rate          float64   `json:"rate"`
	EffectiveDate time.Time `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TaxCalculationRequest struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	State     string  `json:"state"`
}

type TaxCalculationResponse struct {
	Subtotal float64 `json:"subtotal"`
	TaxAmount float64 `json:"tax_amount"`
	Total     float64 `json:"total"`
}

type TaxRateRepository interface {
	GetTaxRate(ctx context.Context, state string, date time.Time) (*TaxRate, error)
	CreateTaxRate(ctx context.Context, rate *TaxRate) error
	UpdateTaxRate(ctx context.Context, rate *TaxRate) error
	DeleteTaxRate(ctx context.Context, id int64) error
}