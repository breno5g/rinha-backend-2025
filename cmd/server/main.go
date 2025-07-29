package main

import (
	"fmt"

	"github.com/breno5g/rinha-back-2025/config"
	"github.com/breno5g/rinha-back-2025/internal/routes"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Init()

	handler := routes.Init()
	logger := config.GetLogger("Main")
	port := fmt.Sprintf(":%d", config.GetEnv().Port)

	logger.Info("Running fasthttp server on port: " + port)
	if err := fasthttp.ListenAndServe(port, handler.HandleRequest); err != nil {
		logger.Errorf("Error in ListenAndServe: %s", err)
	}
}
