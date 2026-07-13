package risk

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) GetRiskPaymentID(w http.ResponseWriter, r *http.Request) {
	paymentIDStr := chi.URLParam(r, "paymentID")
	if paymentIDStr == "" {
		http.Error(w, "paymentID is required", http.StatusBadRequest)
		return
	}

	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		http.Error(w, "invalid paymentID format", http.StatusBadRequest)
		return
	}

	risk, err := h.service.GetRiskPaymentID(paymentID)
	if err != nil {
		http.Error(w, "error retrieving risk: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(risk))
}
