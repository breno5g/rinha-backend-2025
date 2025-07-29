package routes

import (
	"context"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/controller"
	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/health"
	"github.com/breno5g/rinha-back-2025/internal/repository"
	"github.com/breno5g/rinha-back-2025/internal/service"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func Init() *routing.Router {
	db := config.GetDB()
	env := config.GetEnv()
	repo := repository.NewPaymentRepository(db)
	svc := service.NewPaymentService(repo)
	ctrl := controller.NewPaymentController(svc)

	fetcher := &fasthttp.Client{
		ReadTimeout:                   5 * time.Second,
		WriteTimeout:                  5 * time.Second,
		MaxIdleConnDuration:           1 * time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	healthChecker := health.NewHealthChecker(fetcher, env.DefaultURL, env.FallbackURL)
	go healthChecker.Monitor(context.Background(), "default")
	go healthChecker.Monitor(context.Background(), "fallback")

	for i := 0; i < env.MaxWorkers; i++ {
		worker := &entity.Worker{
			Client:        db,
			Repo:          repo,
			WorkerNum:     i,
			Fetcher:       fetcher,
			HealthChecker: healthChecker,
		}
		go worker.Init(context.Background())
	}

	router := routing.New()
	router.Post("/payments", ctrl.Create)
	router.Get("/payments-summary", ctrl.GetSummary)
	router.Get("/health", func(c *routing.Context) error {
		c.SetStatusCode(fasthttp.StatusOK)
		c.SetBodyString("ok")
		return nil
	})

	return router
}
