package app

import (
	"auth-grpc/internal/app/grpcapp"
	"auth-grpc/internal/config"
	"auth-grpc/internal/infrastructure/postgres"
	"auth-grpc/internal/repository"
	"auth-grpc/internal/usecase"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
}

// Init Инициализация зависимостей
func Init(srv *Server, cfg *config.Config, logs *logrus.Logger) *App {
	db, err := postgres.New(cfg)
	if err != nil {
		logs.Error(err) //TODO ...
	}

	repository := repository.New(db)

	uc := usecase.New(logs, repository, repository, cfg.TokenTTL)

	gRPCServer := grpcapp.New(cfg.GRPC.Port, uc, logs)

	return &App{
		GRPCServer: gRPCServer,
	}
}
