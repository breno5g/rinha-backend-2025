package routes

import (
	"net/http"

	"github.com/breno5g/rinha-back-2025/internal/controller"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	controller := controller.NewPaymentController()

	mux.HandleFunc("/payments", controller.Create)
	mux.HandleFunc("/payments-summary", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}
