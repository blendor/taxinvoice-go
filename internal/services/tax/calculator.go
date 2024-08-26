package tax

import (
	"context"
	"fmt"
	"time"

	"github.com/blendor/taxinvoice-go/internal/models"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

type Calculator struct {
	taxRepo     models.TaxRateRepository
	productRepo models.ProductRepository
	logger      *logger.Logger
}

func NewCalculator(tr models.TaxRateRepository, pr models.ProductRepository, l *logger.Logger) *Calculator {
	return &Calculator{
		taxRepo:     tr,
		productRepo: pr,
		logger:      l,
	}
}

func (c *Calculator) Calculate(ctx context.Context, req models.TaxCalculationRequest) (*models.TaxCalculationResponse, error) {
	product, err := c.productRepo.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	taxRate, err := c.taxRepo.GetTaxRate(ctx, req.State, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rate: %w", err)
	}

	subtotal := product.Price * float64(req.Quantity)
	taxAmount := subtotal * taxRate.Rate
	total := subtotal + taxAmount

	return &models.TaxCalculationResponse{
		Subtotal:  subtotal,
		TaxAmount: taxAmount,
		Total:     total,
	}, nil
}