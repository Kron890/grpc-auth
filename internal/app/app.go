package app

import (
	"auth-grpc/internal/app/grpcapp"
	"auth-grpc/internal/config"
	"auth-grpc/internal/infrastructure/storage"
	"auth-grpc/internal/lib/jwt"
	"auth-grpc/internal/repository/postgres"
	"auth-grpc/internal/repository/redisRepo"
	"auth-grpc/internal/usecase"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	RedisServer    *storage.DataBaseRedis
	PostgresServer *storage.DataBase
	GRPCServer     *grpcapp.App
}

// Init Инициализация зависимостей
func Init(srv *Server, cfg *config.Config, logs *logrus.Logger) *App {
	postgresDB, err := storage.NewPostgres(cfg, logs)
	if err != nil {
		logs.Error("database(Postgres) is not connected", err)
	}

	redisDB, err := storage.NewRedis(cfg, logs)
	if err != nil {
		logs.Error("database(Redis) is not connected", err)
	}
	//FIX: вынести в переменную время
	jwtManager, err := jwt.NewManager(15*time.Minute, 7*24*time.Hour)
	if err != nil {
		logs.Error("failed to init jwt manager:", err)
	}

	repoPostgres := postgres.New(postgresDB)

	repoRedis := redisRepo.New(redisDB)

	uc := usecase.New(logs, repoPostgres, repoRedis, jwtManager)

	gRPCServer := grpcapp.New(cfg.GRPC.Port, uc, logs)

	return &App{

		PostgresServer: postgresDB,
		GRPCServer:     gRPCServer,
	}
}
