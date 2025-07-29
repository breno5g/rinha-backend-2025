package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/health"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigFastest

type PaymentRepository interface {
	SaveProcessedPayment(ctx context.Context, payment Payment) error
}

type Worker struct {
	Client        *redis.Client
	Repo          PaymentRepository
	WorkerNum     int
	Fetcher       *fasthttp.Client
	HealthChecker *health.HealthChecker
}

const (
	fastRetryAttempts = 3
	fastRetryDelay    = 100 * time.Millisecond
	requestTimeout    = 5 * time.Second
)

func (w *Worker) Init(ctx context.Context) {
	logger := config.GetLogger(fmt.Sprintf("Worker-%d", w.WorkerNum))
	processingQueue := fmt.Sprintf("payments:processing:%d", w.WorkerNum)
	queue := "payments:queue"

	for {
		result, err := w.Client.BRPopLPush(ctx, queue, processingQueue, 0).Result()
		if err != nil {
			if err != context.Canceled && err != redis.Nil {
				logger.Errorf("Redis BRPopLPush error: %v", err)
				time.Sleep(1 * time.Second)
			}
			continue
		}

		var payment Payment
		if err := json.Unmarshal([]byte(result), &payment); err != nil {
			logger.Errorf("Failed to unmarshal payment: %v", err)
			w.Client.LRem(ctx, processingQueue, 1, result)
			continue
		}

		if !w.processPayment(ctx, payment) {
			logger.Warningf("CRITICAL: Payment %s failed on BOTH processors. Moving to dead-letter queue.", payment.CorrelationId)
			w.Client.LPush(ctx, "payments:dead-letter", result)
		}

		w.Client.LRem(ctx, processingQueue, 1, result)
	}
}

func (w *Worker) processPayment(ctx context.Context, payment Payment) bool {
	env := config.GetEnv()

	payload, _ := json.Marshal(map[string]any{
		"correlationId": payment.CorrelationId.String(),
		"amount":        payment.Amount,
		"requestedAt":   payment.RequestedAt.Format(time.RFC3339Nano),
	})

	var primaryURL, secondaryURL string
	var primaryProcessor, secondaryProcessor string

	if !w.HealthChecker.Default.IsFailing.Load() {
		primaryURL, primaryProcessor = env.DefaultURL, "default"
		secondaryURL, secondaryProcessor = env.FallbackURL, "fallback"
	} else {
		primaryURL, primaryProcessor = env.FallbackURL, "fallback"
		secondaryURL, secondaryProcessor = env.DefaultURL, "default"
	}

	if w.attemptWithRetries(ctx, primaryURL, payload, primaryProcessor, payment) {
		return true
	}

	if w.attemptWithRetries(ctx, secondaryURL, payload, secondaryProcessor, payment) {
		return true
	}

	return false
}

// attemptWithRetries é um helper que encapsula o loop de retentativas rápidas.
func (w *Worker) attemptWithRetries(ctx context.Context, url string, payload []byte, processor string, payment Payment) bool {
	for i := 0; i < fastRetryAttempts; i++ {
		if w.tryProcessor(ctx, url, payload, processor, payment) {
			return true
		}
		time.Sleep(fastRetryDelay)
	}
	return false
}

// tryProcessor realiza uma ÚNICA tentativa de chamada HTTP.
func (w *Worker) tryProcessor(ctx context.Context, url string, payload []byte, processor string, payment Payment) bool {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(payload)

	if err := w.Fetcher.DoTimeout(req, resp, requestTimeout); err != nil {
		return false
	}

	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		payment.Processor = processor
		if err := w.Repo.SaveProcessedPayment(ctx, payment); err != nil {
			config.GetLogger("Worker").Errorf("CRITICAL: Failed to save processed payment %s: %v", payment.CorrelationId, err)
			return true
		}
		return true
	}

	return false
}
