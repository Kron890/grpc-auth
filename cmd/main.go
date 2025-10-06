package main

import (
	"auth-grpc/internal/app"
	"auth-grpc/internal/config"
	"auth-grpc/pgk/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	logs.Info("shutdown signal: ", signalStop)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		application.GRPCServer.StopWithDeadline(ctx)
		application.PostgresServer.Close()
		application.RedisServer.CloseDB()
		close(done)
	}()

	select {
	case <-done:
		logs.Info("shutdown complete")
	case <-ctx.Done():
		logs.Warn("shutdown timed out, forcing exit")
	}

}
