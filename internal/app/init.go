package app

import (
	"auth-grpc/internal/app/grpcapp"
	"auth-grpc/internal/config"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
}

func Init(srv *Server, cfg config.Config, logs *logrus.Logger) *App {
	gRPCServer := grpcapp.New(cfg.GRPC.Port, logs)

	return &App{
		GRPCServer: gRPCServer,
	}
}
