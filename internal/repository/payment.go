package repository

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/breno5g/rinha-back-2025/internal/entity"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

var json = jsoniter.ConfigFastest

const (
	dataHashKey = "payments:data"
)

type PaymentRepository interface {
	AddToQueue(ctx context.Context, payment []byte) error
	SaveProcessedPayment(ctx context.Context, payment *entity.Payment) error
	GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error)
}

type paymentRepository struct {
	db *redis.Client
}

func NewPaymentRepository(db *redis.Client) PaymentRepository {
	return &paymentRepository{db}
}

func (p *paymentRepository) AddToQueue(ctx context.Context, payment []byte) error {
	return p.db.LPush(ctx, "payments:queue", payment).Err()
}

func (p *paymentRepository) SaveProcessedPayment(ctx context.Context, payment *entity.Payment) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	pipe := p.db.Pipeline()
	pipe.HSet(ctx, dataHashKey, payment.CorrelationId.String(), payload)

	timestampIndexKey := "payments:timestamps:" + payment.Processor
	score := float64(payment.RequestedAt.UnixNano())
	pipe.ZAdd(ctx, timestampIndexKey, redis.Z{Score: score, Member: payment.CorrelationId.String()})

	_, err = pipe.Exec(ctx)
	return err
}

func (p *paymentRepository) GetSummary(ctx context.Context, from, to *time.Time) (entity.Summary, error) {
	var summary entity.Summary
	var wg sync.WaitGroup
	var defaultErr, fallbackErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		summary.Default, defaultErr = p.getSummaryForItem(ctx, "default", from, to)
	}()

	go func() {
		defer wg.Done()
		summary.Fallback, fallbackErr = p.getSummaryForItem(ctx, "fallback", from, to)
	}()

	wg.Wait()

	if defaultErr != nil {
		return entity.Summary{}, defaultErr
	}
	if fallbackErr != nil {
		return entity.Summary{}, fallbackErr
	}

	return summary, nil
}

func (p *paymentRepository) getSummaryForItem(ctx context.Context, processor string, from, to *time.Time) (entity.SummaryItem, error) {
	min := "-inf"
	max := "+inf"
	if from != nil {
		min = strconv.FormatInt(from.UnixNano(), 10)
	}
	if to != nil {
		max = strconv.FormatInt(to.UnixNano(), 10)
	}

	timestampIndexKey := "payments:timestamps:" + processor
	paymentIDs, err := p.db.ZRangeByScore(ctx, timestampIndexKey, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()

	if err != nil {
		return entity.SummaryItem{}, err
	}

	if len(paymentIDs) == 0 {
		return entity.SummaryItem{}, nil
	}

	paymentsData, err := p.db.HMGet(ctx, dataHashKey, paymentIDs...).Result()
	if err != nil {
		return entity.SummaryItem{}, err
	}

	var totalAmount float64
	var totalRequests int
	for _, paymentJSON := range paymentsData {
		if paymentJSON == nil {
			continue
		}
		var payment entity.Payment
		if err := json.UnmarshalFromString(paymentJSON.(string), &payment); err == nil {
			totalAmount += payment.Amount
			totalRequests++
		}
	}

	return entity.SummaryItem{
		TotalRequests: totalRequests,
		TotalAmount:   totalAmount,
	}, nil
}
