package service

import (
	"context"
	"fmt"

	"github.com/blendor/taxinvoice-go/internal/model"
	"github.com/blendor/taxinvoice-go/internal/store"
)

type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) CalculateTax(ctx context.Context, req model.TaxRequest) (*model.TaxResponse, error) {
	product, err := s.store.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	taxRate, err := s.store.GetTaxRate(ctx, req.State)
	if err != nil {
		return nil, fmt.Errorf("tax rate not found: %w", err)
	}

	subtotal := product.Price * float64(req.Quantity)
	taxAmount := subtotal * (taxRate.Rate / 100)
	total := subtotal + taxAmount

	return &model.TaxResponse{
		Subtotal:  subtotal,
		TaxAmount: taxAmount,
		Total:     total,
	}, nil
}

func (s *Service) CreateInvoice(ctx context.Context, req model.InvoiceRequest) (*model.Invoice, error) {
	taxRate, err := s.store.GetTaxRate(ctx, req.State)
	if err != nil {
		return nil, fmt.Errorf("tax rate not found: %w", err)
	}

	var items []model.InvoiceItem
	var subtotal float64

	for _, item := range req.Items {
		product, err := s.store.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", item.ProductID, err)
		}

		itemSubtotal := product.Price * float64(item.Quantity)
		subtotal += itemSubtotal

		items = append(items, model.InvoiceItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
			Subtotal:  itemSubtotal,
		})
	}

	taxAmount := subtotal * (taxRate.Rate / 100)
	total := subtotal + taxAmount

	invoice := &model.Invoice{
		CustomerID: req.CustomerID,
		State:      req.State,
		Subtotal:   subtotal,
		TaxAmount:  taxAmount,
		Total:      total,
		Items:      items,
	}

	if err := s.store.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}
