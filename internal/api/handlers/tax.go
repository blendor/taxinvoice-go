package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blendor/taxinvoice-go/internal/models"
	"github.com/blendor/taxinvoice-go/internal/services/tax"
	"github.com/blendor/taxinvoice-go/pkg/logger"
)

type TaxHandler struct {
	calculator *tax.Calculator
	logger     *logger.Logger
}

func NewTaxHandler(calculator *tax.Calculator, logger *logger.Logger) *TaxHandler {
	return &TaxHandler{
		calculator: calculator,
		logger:     logger,
	}
}

func (h *TaxHandler) CalculateTax(w http.ResponseWriter, r *http.Request) {
	var request models.TaxCalculationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.calculator.Calculate(r.Context(), request)
	if err != nil {
		h.logger.Error("Failed to calculate tax", "error", err)
		http.Error(w, "Failed to calculate tax", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}