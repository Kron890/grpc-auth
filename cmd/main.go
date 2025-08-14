package main

import (
	"auth-grpc/internal/app"
	"auth-grpc/internal/config"
	"auth-grpc/pgk/logger"

	"github.com/labstack/gommon/log"
)

func main() {
	cfg := config.MustLoad()

	srv := app.NewServer()

	logs := logger.Init()
	logs.Info("logger is initialized")
	logs.Info(cfg)

	err := app.Init(srv, *cfg, logs)
	if err != nil {
		log.Error("Error: failed to initialize application:", err)
	}

}
