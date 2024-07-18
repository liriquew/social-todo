package grpcauth

import (
	"context"
	"errors"
	"fmt"
	"sso_service/internal/sevices/auth"

	"github.com/liriquew/todoprotos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(context.Context, string, string) (string, error)
	Register(context.Context, string, string) (int64, error)
	Authorize(context.Context, string) (int64, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (g *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if err := validateRequest(req.Username, req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := g.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "invalid username or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &sso.LoginResponse{Token: token}, nil
}

func (g *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if err := validateRequest(req.Username, req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	uid, err := g.auth.Register(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExist) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.RegisterResponse{Uid: uid}, nil
}

func (g *serverAPI) Authorize(ctx context.Context, req *sso.AuthorizeRequest) (*sso.AuthorizeResponse, error) {
	token := req.Token
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "jwt token is empty")
	}

	uid, err := g.auth.Authorize(ctx, token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &sso.AuthorizeResponse{Uid: uid}, nil
}

func validateRequest(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
