package friends_grpc

import (
	"context"
	"fmt"

	"github.com/liriquew/todoprotos/gen/go/friends"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendsService interface {
	AddFriend(context.Context, int64, int64) error
	RemoveFriend(context.Context, int64, int64) error
	ListUserFriends(context.Context, int64) ([]int64, error)
}

type serverAPI struct {
	friends.UnimplementedFriendsServer
	api FriendsService
}

func Register(gRPC *grpc.Server, friendsService FriendsService) {
	friends.RegisterFriendsServer(gRPC, &serverAPI{api: friendsService})
}

var (
	ErrBadUID      = fmt.Errorf("bad UID")
	ErrBadFriendID = fmt.Errorf("bad friend ID")
)

func (s *serverAPI) AddFriend(ctx context.Context, req *friends.FriendRequest) (*friends.FriendResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.AddFriend(ctx, req.UID, req.FriendID)
	if err != nil {
		// TODO: check err

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &friends.FriendResponse{}, nil
}

func (s *serverAPI) RemoveFriend(ctx context.Context, req *friends.FriendRequest) (*friends.FriendResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.api.RemoveFriend(ctx, req.UID, req.FriendID)
	if err != nil {
		// TODO: check err

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &friends.FriendResponse{}, nil
}

func (s *serverAPI) ListFriends(ctx context.Context, req *friends.ListFriendRequest) (*friends.ListFriendResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	friendIDs, err := s.api.ListUserFriends(ctx, req.UID)
	if err != nil {
		// TODO: check err

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &friends.ListFriendResponse{FriendIDs: friendIDs}, nil
}

func validateRequest(req interface{}) error {
	switch val := req.(type) {
	case *friends.FriendRequest:
		if val.UID <= 0 {
			return ErrBadUID
		}
		if val.FriendID <= 0 {
			return ErrBadFriendID
		}
	case *friends.ListFriendRequest:
		if val.UID <= 0 {
			return ErrBadUID
		}
	}
	return nil
}
