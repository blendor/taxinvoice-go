package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blendor/taxinvoice-go/internal/models"
	"github.com/blendor/taxinvoice-go/internal/services/invoice"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

type InvoiceHandler struct {
	generator *invoice.Generator
	logger    *logger.Logger
}

func NewInvoiceHandler(generator *invoice.Generator, logger *logger.Logger) *InvoiceHandler {
	return &InvoiceHandler{
		generator: generator,
		logger:    logger,
	}
}

func (h *InvoiceHandler) GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	var request models.InvoiceGenerationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.generator.Generate(r.Context(), request)
	if err != nil {
		h.logger.Error("Failed to generate invoice", "error", err)
		http.Error(w, "Failed to generate invoice", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}