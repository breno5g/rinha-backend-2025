package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/service"
	"github.com/breno5g/rinha-back-2025/internal/utils"
	"github.com/google/uuid"
)

type PaymentController struct {
	svc service.PaymentService
}

func NewPaymentController(svc service.PaymentService) *PaymentController {
	return &PaymentController{
		svc,
	}
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

	if req.CorrelationId == uuid.Nil {
		logger.Errorf("Invalid correlationId: %v", req.CorrelationId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		logger.Errorf("Invalid amount: %v", req.Amount)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := p.svc.AddToQueue(r.Context(), req)
	if err != nil {
		logger.Errorf("Failed to add payment to queue: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *PaymentController) GetSummary(w http.ResponseWriter, r *http.Request) {
	logger := config.GetLogger("Get Payment Summary")

	from := utils.ParseTime(r.URL.Query().Get("from"))
	to := utils.ParseTime(r.URL.Query().Get("to"))

	summary, err := p.svc.GetSummary(r.Context(), from, to)
	if err != nil {
		logger.Errorf("Failed to get summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
