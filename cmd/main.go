package main

import (
	"auth-grpc/internal/app"
	"auth-grpc/internal/config"
	"auth-grpc/pgk/logger"
)

func main() {
	cfg := config.MustLoad()

	srv := app.NewServer()

	logs := logger.Init()
	logs.Info("logger is initialized")

	initializer := app.Init(srv, *cfg, logs)

	initializer.GRPCServer.MustRun()

}
