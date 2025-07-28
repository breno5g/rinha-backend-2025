package main

import (
	"net/http"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/routes"
)

func main() {
	config.Init()

	mux := routes.Init()
	logger := config.GetLogger("Main")
	logger.Info("Running server")
	http.ListenAndServe(":8080", mux)
}
