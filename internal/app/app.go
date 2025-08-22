package app

import (
	"auth-grpc/internal/app/grpcapp"
	"auth-grpc/internal/config"
	"auth-grpc/internal/infrastructure/storage"
	"auth-grpc/internal/repository/postgres"
	"auth-grpc/internal/usecase"

	"github.com/sirupsen/logrus"
)

type App struct {
	PostgresServer *storage.DataBase
	GRPCServer     *grpcapp.App
}

// Init Инициализация зависимостей
func Init(srv *Server, cfg *config.Config, logs *logrus.Logger) *App {
	PostgresDB, err := storage.NewPostgres(cfg)
	if err != nil {
		logs.Error("database is not connected", err)
	}

	repository := postgres.New(PostgresDB)

	uc := usecase.New(logs, repository, repository, cfg.TokenTTL)

	gRPCServer := grpcapp.New(cfg.GRPC.Port, uc, logs)

	return &App{
		PostgresServer: PostgresDB,
		GRPCServer:     gRPCServer,
	}
}
