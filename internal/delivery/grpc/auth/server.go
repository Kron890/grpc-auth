package auth

import (
	sso "auth-grpc/contract/gen/auth"
	"context"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, login, password string) (string, error)
	login(ctx context.Context, login, password string) error
}

type ServerAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

// Регистрируеми сервер
func RegisterServer(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

// Регистрация пользователя
func (s *ServerAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {

	if !govalidator.IsASCII(req.GetLogin()) {
		return &sso.RegisterResponse{}, status.Error(codes.InvalidArgument, "email is required")
	}

	if !govalidator.IsASCII(req.GetPassword()) {
		return &sso.RegisterResponse{}, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		//todo...
		return &sso.RegisterResponse{}, status.Error(codes.Internal, "internal error")
	}

	return &sso.RegisterResponse{UserId: id}, nil
}

// Авторизация пользователя
func (s *ServerAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if !govalidator.IsASCII(req.GetLogin()) {
		return &sso.LoginResponse{}, status.Error(codes.InvalidArgument, "email is required")
	}

	if !govalidator.IsASCII(req.GetPassword()) {
		return &sso.LoginResponse{}, status.Error(codes.InvalidArgument, "incorrect password")
	}

	err := s.auth.login(ctx, req.GetLogin(), req.Password)
	if err != nil {
		//todo ...
		return &sso.LoginResponse{}, status.Error(codes.Internal, "internal error")
	}
	return &sso.LoginResponse{}, nil
}

func (s *ServerAPI) VerifyToken(ctx context.Context, req *sso.VerifyTokenRequest) (*sso.VerifyTokenResponse, error) {
	panic("VerifyToken")
}

func (s *ServerAPI) RefreshTokens(ctx context.Context, req *sso.RefreshTokensRequest) (*sso.RefreshTokensResponse, error) {
	panic("RefreshTokens")
}
