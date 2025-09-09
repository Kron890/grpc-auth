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

	application := app.Init(srv, cfg, logs)

	go application.GRPCServer.MustRun()

	// shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	signalStop := <-stop
	application.PostgresServer.Close()
	application.RedisServer.CloseDB()
	application.GRPCServer.Stop()

	logs.Info("initializer stopped: ", signalStop)

}
