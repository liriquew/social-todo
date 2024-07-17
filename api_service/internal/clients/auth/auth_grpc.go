package auth_grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/liriquew/todoprotos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Client struct {
	api sso.AuthClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "auth_grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
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

// TODO: реализовать методы для клиента
