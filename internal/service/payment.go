package service

import (
	"context"
	"encoding/json"

	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/repository"
)

type PaymentService interface {
	AddToQueue(ctx context.Context, payment entity.Payment) error
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
