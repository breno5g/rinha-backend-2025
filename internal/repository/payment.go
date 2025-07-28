package repository

import (
	"context"
	"encoding/json"

	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/redis/go-redis/v9"
)

type PaymentRepository interface {
	AddToQueue(ctx context.Context, payment []byte) error
}

type paymentRepository struct {
	db *redis.Client
}

func NewPaymentRepository(db *redis.Client) *paymentRepository {
	return &paymentRepository{
		db,
	}
}

func (p *paymentRepository) AddToQueue(ctx context.Context, payment []byte) error {
	err := p.db.LPush(ctx, "payments:queue", payment).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *paymentRepository) Save(ctx context.Context, payment entity.Payment) error {
	json, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	return p.db.HSet(ctx, "payments", payment.CorrelationId, json).Err()
}

func (p *paymentRepository) GetAll(ctx context.Context) ([]entity.Payment, error) {
	paymentsData, err := p.db.HGetAll(ctx, "payments").Result()
	if err != nil {
		return nil, err
	}
	var payments []entity.Payment
	for _, paymentDataJSON := range paymentsData {
		var paymentData entity.Payment
		if err := json.Unmarshal([]byte(paymentDataJSON), &paymentData); err != nil {
			continue
		}

		payments = append(payments, paymentData)
	}
	return payments, nil
}
