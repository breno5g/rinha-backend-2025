package main

import (
	"fmt"
	"net/http"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/routes"
)

func main() {
	config.Init()

	mux := routes.Init()
	logger := config.GetLogger("Main")
	port := fmt.Sprintf(":%d", config.GetEnv().Port)
	logger.Info("Running server on port: " + port)
	http.ListenAndServe(port, mux)
}
