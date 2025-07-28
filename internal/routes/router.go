package routes

import (
	"net/http"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/controller"
	"github.com/breno5g/rinha-back-2025/internal/repository"
	"github.com/breno5g/rinha-back-2025/internal/service"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	db := config.GetDB()
	repository := repository.NewPaymentRepository(db)
	service := service.NewPaymentService(repository)
	controller := controller.NewPaymentController(service)

	mux.HandleFunc("/payments", controller.Create)
	mux.HandleFunc("/payments-summary", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}
