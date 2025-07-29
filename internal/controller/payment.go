package controller

import (
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/service"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigFastest

type PaymentController struct {
	svc service.PaymentService
}

func NewPaymentController(svc service.PaymentService) *PaymentController {
	return &PaymentController{
		svc: svc,
	}
}

// parseTime foi movido para cá como uma função auxiliar privada.
func parseTime(param []byte) *time.Time {
	if len(param) == 0 {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, string(param)); err == nil {
		return &t
	}
	return nil
}

func (p *PaymentController) Create(c *routing.Context) error {
	var req entity.Payment
	if err := json.Unmarshal(c.PostBody(), &req); err != nil {
		c.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	req.RequestedAt = time.Now().UTC()

	if req.CorrelationId == uuid.Nil || req.Amount <= 0 {
		c.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	err := p.svc.AddToQueue(c, req)
	if err != nil {
		config.GetLogger("CreatePayment").Errorf("Failed to add payment to queue: %v", err)
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil
	}

	c.SetStatusCode(fasthttp.StatusAccepted)
	return nil
}

func (p *PaymentController) GetSummary(c *routing.Context) error {
	from := parseTime(c.QueryArgs().Peek("from"))
	to := parseTime(c.QueryArgs().Peek("to"))

	summary, err := p.svc.GetSummary(c, from, to)
	if err != nil {
		config.GetLogger("GetSummary").Errorf("Failed to get summary: %v", err)
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil
	}

	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(c).Encode(summary)
	return nil
}
