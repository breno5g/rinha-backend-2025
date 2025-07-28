package main

import (
	"net/http"

	"github.com/breno5g/rinha-back-2025/internal/routes"
)

func main() {
	// config.Init()

	mux := routes.Init()
	http.ListenAndServe(":8080", mux)
}
