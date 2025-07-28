package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/google/uuid"
)

type PaymentController struct{}

func NewPaymentController() *PaymentController {
	return &PaymentController{}
}

func (p *PaymentController) Create(w http.ResponseWriter, r *http.Request) {
	logger := config.GetLogger("Create Payment")

	req := entity.Payment{
		RequestedAt: time.Now().UTC(),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("Invalid body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(req.CorrelationId)
	if err != nil {
		logger.Errorf("Invalid correlationId: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		logger.Errorf("Invalid amount: %v", req.Amount)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(req)
}
