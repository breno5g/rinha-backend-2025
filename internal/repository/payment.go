package repository

import (
	"context"

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
