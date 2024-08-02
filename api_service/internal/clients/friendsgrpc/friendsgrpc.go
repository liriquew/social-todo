package friends_grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/todoprotos/gen/go/friends"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Client struct {
	api friends.FriendsClient
	log *slog.Logger
}

func New(log *slog.Logger, cfg config.ServiceConfig) (*Client, error) {
	const op = "friends_grpc.New"

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
		api: friends.NewFriendsClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

var (
	ErrNotFound        = fmt.Errorf("note not found")
	ErrAlreadyExists   = fmt.Errorf("note already exists")
	ErrInvalidArgument = fmt.Errorf("invalid argument")
)

func (c *Client) AddFriend(ctx context.Context, UID, FID int64) error {
	const op = "friends_grpc.Create"

	_, err := c.api.AddFriend(ctx, &friends.FriendRequest{UID: UID, FriendID: FID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) RemoveFriend(ctx context.Context, UID, FID int64) error {
	const op = "friends_grpc.Get"

	_, err := c.api.RemoveFriend(ctx, &friends.FriendRequest{
		UID:      UID,
		FriendID: FID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) ListFriends(ctx context.Context, UID int64) ([]int64, error) {
	const op = "friends_grpc.Update"

	resp, err := c.api.ListFriends(ctx, &friends.ListFriendRequest{
		UID: UID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp.FriendIDs, nil
}
