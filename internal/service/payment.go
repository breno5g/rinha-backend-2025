package service

import (
	"context"
	"time"

	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/repository"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type PaymentService interface {
	AddToQueue(ctx context.Context, payment entity.Payment) error
	GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

func (p *paymentService) AddToQueue(ctx context.Context, payment entity.Payment) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return p.repo.AddToQueue(ctx, data)
}

func (p *paymentService) GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error) {
	return p.repo.GetSummary(ctx, from, to)
}
