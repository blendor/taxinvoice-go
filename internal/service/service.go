package service

import (
	"context"
	"fmt"
	"time"

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

func (s *Service) GetInvoice(ctx context.Context, id int64) (*model.Invoice, error) {
	inv, err := s.store.GetInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}
	s.checkOverdue(inv)
	return inv, nil
}

func (s *Service) ListInvoices(ctx context.Context, status string) ([]model.Invoice, error) {
	invoices, err := s.store.ListInvoices(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	for i := range invoices {
		s.checkOverdue(&invoices[i])
	}
	return invoices, nil
}

func (s *Service) MarkPaid(ctx context.Context, id int64) (*model.Invoice, error) {
	inv, err := s.store.GetInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}
	if inv.Status == model.StatusPaid {
		return inv, nil
	}
	if inv.Status == model.StatusRefunded {
		return nil, fmt.Errorf("cannot mark refunded invoice as paid")
	}
	now := time.Now()
	if err := s.store.UpdateInvoiceStatus(ctx, id, model.StatusPaid, &now); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}
	inv.Status = model.StatusPaid
	inv.PaidAt = &now
	return inv, nil
}

func (s *Service) MarkRefunded(ctx context.Context, id int64) (*model.Invoice, error) {
	inv, err := s.store.GetInvoice(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}
	if inv.Status != model.StatusPaid {
		return nil, fmt.Errorf("can only refund paid invoices")
	}
	if err := s.store.UpdateInvoiceStatus(ctx, id, model.StatusRefunded, inv.PaidAt); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}
	inv.Status = model.StatusRefunded
	return inv, nil
}

func (s *Service) checkOverdue(inv *model.Invoice) {
	if inv.Status == model.StatusUnpaid && time.Now().After(inv.DueDate) {
		inv.Status = model.StatusOverdue
	}
}
