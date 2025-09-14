package auth

import (
	sso "auth-grpc/contract/gen/auth"
	"auth-grpc/internal"
	"context"
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	sso.UnimplementedAuthServer
	auth internal.Auth
}

// New регистрируеми сервер
func New(gRPC *grpc.Server, auth internal.Auth) {
	sso.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

// Register регистрация пользователя
func (s *ServerAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if err := s.checkLoginPass(req.GetLogin(), req.Password); err != nil {
		return &sso.RegisterResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.auth.Register(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, internal.ErrInvailidCredentials) {
			return &sso.RegisterResponse{}, status.Error(codes.Internal, "invalid credentials")
		}
		return &sso.RegisterResponse{}, status.Error(codes.Internal, "internal error")
	}

	return &sso.RegisterResponse{UserId: fmt.Sprintf("%d", id)}, nil
}

// Login авторизация пользователя
func (s *ServerAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if err := s.checkLoginPass(req.GetLogin(), req.Password); err != nil {
		return &sso.LoginResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetLogin(), req.Password)
	if err != nil {
		if errors.Is(err, internal.ErrInvailidCredentials) {
			return &sso.LoginResponse{}, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return &sso.LoginResponse{}, status.Error(codes.Internal, "internal error")
	}
	return &sso.LoginResponse{AccessToken: token.Access}, nil
}

func (s *ServerAPI) VerifyToken(ctx context.Context, req *sso.VerifyTokenRequest) (*sso.VerifyTokenResponse, error) {
	panic("VerifyToken")
}

func (s *ServerAPI) RefreshTokens(ctx context.Context, req *sso.RefreshTokensRequest) (*sso.RefreshTokensResponse, error) {
	panic("RefreshTokens")
}

func (s *ServerAPI) checkLoginPass(login string, password string) error {
	if !govalidator.IsASCII(login) {
		return fmt.Errorf("email is required")
	}

	if !govalidator.IsASCII(password) {
		return fmt.Errorf("password is required")
	}
	return nil
}
