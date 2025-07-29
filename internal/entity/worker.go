package entity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/redis/go-redis/v9"
)

type PaymentRepository interface {
	GetAll(ctx context.Context) ([]Payment, error)
	Save(ctx context.Context, payment Payment) error
}

type Worker struct {
	Client    *redis.Client
	Repo      PaymentRepository
	WorkerNum int
	Fetcher   *http.Client
}

const (
	retryDelay    = time.Second
	maxRetries    = 5
	retryInterval = 500 * time.Millisecond
	httpTimeout   = 5 * time.Second
)

func (w *Worker) Init(ctx context.Context) {
	logger := config.GetLogger("Payment workers")
	processingQueue := fmt.Sprintf("payments:processing:%d", w.WorkerNum)
	w.Fetcher = &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * retryDelay,
			DisableKeepAlives:   false,
		},
	}

	for {
		result, err := w.Client.RPopLPush(ctx, "payments:queue", processingQueue).Result()
		if err != nil {
			if err == redis.Nil {
				time.Sleep(1 * retryDelay)
				continue
			}
			logger.Debugf("[Worker %d] Redis error: %v", w.WorkerNum, err)
			time.Sleep(1 * retryDelay)
			continue
		}
		var payment Payment
		if err := json.Unmarshal([]byte(result), &payment); err != nil {
			logger.Debugf("[Worker %d] Failed to unmarshal payment: %v", w.WorkerNum, err)
			w.Client.LRem(ctx, processingQueue, 1, result)
			continue
		}
		if !w.processPayment(ctx, payment) {
			w.Client.LPush(ctx, "payments:queue", result)
			w.Client.LRem(ctx, processingQueue, 1, result)
			continue
		}
		w.Client.LRem(ctx, processingQueue, 1, result)
	}
}

func (w *Worker) processPayment(ctx context.Context, payment Payment) bool {
	defaultURL := config.GetEnv().DefaultURL
	fallbackURL := config.GetEnv().FallbackURL

	payload, _ := json.Marshal(map[string]any{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	})

	for range maxRetries {
		if w.tryProcessor(ctx, defaultURL, payload, "default", payment) {
			return true
		}
		time.Sleep(retryInterval)
	}

	return w.tryProcessor(ctx, fallbackURL, payload, "fallback", payment)
}

func (w *Worker) tryProcessor(ctx context.Context, url string, payload []byte, processor string, payment Payment) bool {
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.Fetcher.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		payment.Processor = processor
		w.Repo.Save(ctx, payment)
		return true
	}
	return false
}
