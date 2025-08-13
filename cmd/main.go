package main

import (
	"auth-grpc/config"
	"auth-grpc/internal/app"
	"auth-grpc/logger"
)

func main() {
	logger.Init()

	srv := app.NewServer()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Log.Error("Error: unable to load configuration:", err)
	}

	cfgGRPC, err := config.GetConfigGRPC()
	if err != nil {
		logger.Log.Error("Error: unable to load configuration:", err)
	}

	err = app.Init(srv, cfg, cfgGRPC)
	if err != nil {
		logger.Log.Error("Error: failed to initialize application:", err)
	}

}
