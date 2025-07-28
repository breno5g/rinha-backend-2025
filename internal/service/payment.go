package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/repository"
)

type PaymentService interface {
	AddToQueue(ctx context.Context, payment entity.Payment) error
	GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *paymentService {
	return &paymentService{
		repo,
	}
}

func (p *paymentService) AddToQueue(ctx context.Context, payment entity.Payment) error {
	data, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	err = p.repo.AddToQueue(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (p *paymentService) GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error) {
	payments, err := p.repo.GetAll(ctx)
	if err != nil {
		return entity.Summary{}, err
	}

	var summary entity.Summary
	for _, payment := range payments {
		if (from != nil && !from.IsZero() && payment.RequestedAt.Before(*from)) ||
			(to != nil && !to.IsZero() && payment.RequestedAt.After(*to)) {
			continue
		}

		item := &summary.Default
		if payment.Processor == "fallback" {
			item = &summary.Fallback
		}
		if payment.Processor == "default" || payment.Processor == "fallback" {
			item.TotalRequests++
			item.TotalAmount += payment.Amount
		}
	}
	return summary, nil
}
