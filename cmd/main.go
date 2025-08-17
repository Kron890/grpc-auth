package main

import (
	"auth-grpc/internal/app"
	"auth-grpc/internal/config"
	"auth-grpc/pgk/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	srv := app.NewServer()

	logs := logger.Init()
	logs.Info("logger is initialized")

	initializer := app.Init(srv, *cfg, logs)

	go initializer.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	signalStop := <-stop

	initializer.GRPCServer.Stop()
	logs.Info("initializer stopped: ", signalStop)

}
