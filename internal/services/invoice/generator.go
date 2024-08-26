package invoice

import (
	"context"
	"fmt"

	"github.com/yourusername/tax-invoice-app/internal/models"
	"github.com/yourusername/tax-invoice-app/pkg/logger"
)

type Generator struct {
	invoiceRepo models.InvoiceRepository
	taxCalc     *tax.Calculator
	pdfGen      *PDFGenerator
	csvGen      *CSVGenerator
	logger      *logger.Logger
}

func NewGenerator(ir models.InvoiceRepository, tc *tax.Calculator, pg *PDFGenerator, cg *CSVGenerator, l *logger.Logger) *Generator {
	return &Generator{
		invoiceRepo: ir,
		taxCalc:     tc,
		pdfGen:      pg,
		csvGen:      cg,
		logger:      l,
	}
}

func (g *Generator) Generate(ctx context.Context, req models.InvoiceGenerationRequest) (*models.InvoiceGenerationResponse, error) {
	var totalAmount, taxAmount float64

	for i, item := range req.Items {
		product, err := g.taxCalc.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product for item %d: %w", i, err)
		}

		taxReq := models.TaxCalculationRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			State:     req.State,
		}

		taxResp, err := g.taxCalc.Calculate(ctx, taxReq)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax for item %d: %w", i, err)
		}

		totalAmount += taxResp.Total
		taxAmount += taxResp.TaxAmount
		req.Items[i].UnitPrice = product.Price
		req.Items[i].Subtotal = taxResp.Subtotal
	}

	invoice := &models.Invoice{
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		TaxAmount:   taxAmount,
		Items:       req.Items,
	}

	if err := g.invoiceRepo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	pdfLocation, err := g.pdfGen.Generate(invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	csvLocation, err := g.csvGen.Generate(invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CSV: %w", err)
	}

	return &models.InvoiceGenerationResponse{
		Invoice:     *invoice,
		PDFLocation: pdfLocation,
		CSVLocation: csvLocation,
	}, nil
}