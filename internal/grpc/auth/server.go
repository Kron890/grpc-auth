package auth

import (
	"auth-grpc/contract/gen/auth"
	"context"

	"google.golang.org/grpc"
)

type ServerAPI struct {
	auth.UnimplementedAuthServer
}

func RegisterServer(gRPC *grpc.Server) {
	auth.RegisterAuthServer(gRPC, &ServerAPI{})
}
func Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	panic("register")
}

func Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	panic("login")
}

func VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	panic("VerifyToken")
}

func RefreshTokens(ctx context.Context, req *auth.RefreshTokensRequest) (*auth.RefreshTokensResponse, error) {
	panic(RefreshTokens)
}
