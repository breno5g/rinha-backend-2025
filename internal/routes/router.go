package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/controller"
	"github.com/breno5g/rinha-back-2025/internal/entity"
	"github.com/breno5g/rinha-back-2025/internal/repository"
	"github.com/breno5g/rinha-back-2025/internal/service"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	db := config.GetDB()
	repository := repository.NewPaymentRepository(db)
	service := service.NewPaymentService(repository)
	controller := controller.NewPaymentController(service)

	fetcher := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
		},
	}

	for i := 0; i < config.GetEnv().MaxWorkers; i++ {
		worker := &entity.Worker{Client: db, Repo: repository, WorkerNum: i, Fetcher: fetcher}
		go worker.Init(context.Background())
	}

	mux.HandleFunc("/payments", controller.Create)
	mux.HandleFunc("/payments-summary", controller.GetSummary)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}
