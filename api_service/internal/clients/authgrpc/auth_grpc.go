package auth_grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/todoprotos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Client struct {
	api sso.AuthClient
	log *slog.Logger
}

var (
	ErrMissJWTToken = fmt.Errorf("miss JWT auth token")
)

func New(log *slog.Logger, cfg config.ServiceConfig) (*Client, error) {
	const op = "auth_grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.Retries)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(cfg.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: sso.NewAuthClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

var (
	ErrNotFound        = fmt.Errorf("user not found")
	ErrAlreadyExists   = fmt.Errorf("user already exists")
	ErrInvalidArgument = fmt.Errorf("invalid argument")
)

func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	const op = "auth_grpc.Login"

	resp, err := c.api.Login(ctx, &sso.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return "", ErrNotFound
			case codes.InvalidArgument:
				return "", ErrInvalidArgument
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}

func (c *Client) Register(ctx context.Context, username, password string) (int64, error) {
	const op = "auth_grpc.Login"

	resp, err := c.api.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return 0, ErrAlreadyExists
			case codes.InvalidArgument:
				return 0, ErrInvalidArgument
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Uid, nil
}

func (c *Client) Authorize(ctx context.Context, token string) (int64, error) {
	const op = "auth_grpc.Authorize"

	resp, err := c.api.Authorize(ctx, &sso.AuthorizeRequest{
		Token: token,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return 0, ErrMissJWTToken
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Uid, nil
}
