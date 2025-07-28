package routes

import (
	"net/http"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/payments-summary", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}
