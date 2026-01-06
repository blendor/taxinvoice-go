package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/blendor/taxinvoice-go/internal/model"
	"github.com/blendor/taxinvoice-go/internal/service"
)

type Handler struct {
	svc *service.Service
	log *slog.Logger
}

func New(svc *service.Service, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) CalculateTax(w http.ResponseWriter, r *http.Request) {
	var req model.TaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProductID == 0 || req.Quantity <= 0 || len(req.State) != 2 {
		h.error(w, "invalid request: product_id, quantity > 0, and 2-letter state required", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CalculateTax(r.Context(), req)
	if err != nil {
		h.log.Error("calculate tax failed", "error", err)
		h.error(w, "failed to calculate tax", http.StatusInternalServerError)
		return
	}

	h.json(w, resp)
}

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req model.InvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.CustomerID == 0 || len(req.State) != 2 || len(req.Items) == 0 {
		h.error(w, "invalid request: customer_id, 2-letter state, and items required", http.StatusBadRequest)
		return
	}

	invoice, err := h.svc.CreateInvoice(r.Context(), req)
	if err != nil {
		h.log.Error("create invoice failed", "error", err)
		h.error(w, "failed to create invoice", http.StatusInternalServerError)
		return
	}

	h.json(w, invoice)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.json(w, map[string]string{"status": "ok"})
}

func (h *Handler) json(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) error(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
