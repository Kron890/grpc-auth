package app

import (
	"auth-grpc/infrastructure/postgres"
	"auth-grpc/internal/app/grpcapp"
	"auth-grpc/internal/config"
	"auth-grpc/internal/repository"
	"auth-grpc/internal/usecase"
	"fmt"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *grpcapp.App
}

// Init Инициализация зависимостей
func Init(srv *Server, cfg config.Config, logs *logrus.Logger) *App {
	db, err := postgres.NewDBConnect(cfg.DBPort)
	if err != nil {
		logs.Error(err) //TODO ...
	}

	repository, err := repository.New(db)
	uc := usecase.New(*logs, repository, repository, cfg.TokenTTL)
	fmt.Println(uc)

	gRPCServer := grpcapp.New(cfg.GRPC.Port, logs)

	return &App{
		GRPCServer: gRPCServer,
	}
}
