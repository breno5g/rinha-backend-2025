package health

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigFastest

type HealthStatus struct {
	IsFailing atomic.Bool
}

type HealthChecker struct {
	client      *fasthttp.Client
	defaultURL  string
	fallbackURL string
	Default     HealthStatus
	Fallback    HealthStatus
}

type healthCheckResponse struct {
	Failing bool `json:"failing"`
}

func NewHealthChecker(client *fasthttp.Client, defaultURL, fallbackURL string) *HealthChecker {
	return &HealthChecker{
		client:      client,
		defaultURL:  defaultURL + "/service-health",
		fallbackURL: fallbackURL + "/service-health",
	}
}

func (h *HealthChecker) Monitor(ctx context.Context, processor string) {
	url := h.defaultURL
	status := &h.Default
	if processor == "fallback" {
		url = h.fallbackURL
		status = &h.Fallback
	}

	logger := config.GetLogger("HealthChecker-" + processor)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			req.SetRequestURI(url)
			req.Header.SetMethod(fasthttp.MethodGet)

			err := h.client.DoTimeout(req, resp, 2*time.Second)
			if err != nil {
				logger.Errorf("Health check failed: %v", err)
				status.IsFailing.Store(true)
			} else if resp.StatusCode() != fasthttp.StatusOK {
				logger.Warningf("Health check returned status %d", resp.StatusCode())
				status.IsFailing.Store(true)
			} else {
				var healthResp healthCheckResponse
				if jsonErr := json.Unmarshal(resp.Body(), &healthResp); jsonErr != nil {
					logger.Errorf("Failed to unmarshal health response: %v", jsonErr)
					status.IsFailing.Store(true)
				} else {
					status.IsFailing.Store(healthResp.Failing)
				}
			}

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}
	}
}
