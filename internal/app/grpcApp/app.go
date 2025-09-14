package grpcapp

import (
	authgrpc "auth-grpc/internal/delivery/grpc/auth"
	"auth-grpc/internal/delivery/middleware"
	"auth-grpc/internal/jwt"
	"auth-grpc/internal/usecase"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	logs       *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(port int, uc *usecase.Auth, jwtManager jwt.JWTManager, logs *logrus.Logger) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.NewAuthInterceptor(jwtManager)))

	authgrpc.New(gRPCServer, uc)

	return &App{
		logs:       logs,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Запуск gRPC сервера
// TODO: Убрать must
func (a *App) MustRun() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err)
	}
	a.logs.Info("gRPC server is running", listen.Addr().String())

	if err := a.gRPCServer.Serve(listen); err != nil {
		panic(err)
	}

}

func (a *App) Stop() {
	a.logs.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()

}
