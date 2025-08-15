package auth

import (
	sso "auth-grpc/contract/gen/auth"
	"context"

	"google.golang.org/grpc"
)

type ServerAPI struct {
	sso.UnimplementedAuthServer
}

func RegisterServer(gRPC *grpc.Server) {
	sso.RegisterAuthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	panic("register")
}

func (s *ServerAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	panic("login")
}

func (s *ServerAPI) VerifyToken(ctx context.Context, req *sso.VerifyTokenRequest) (*sso.VerifyTokenResponse, error) {
	panic("VerifyToken")
}

func (s *ServerAPI) RefreshTokens(ctx context.Context, req *sso.RefreshTokensRequest) (*sso.RefreshTokensResponse, error) {
	panic("RefreshTokens")
}
