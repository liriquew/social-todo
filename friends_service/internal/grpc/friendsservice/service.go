package friendssrvc

import (
	"context"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
	neo_storage "github.com/liriquew/social-todo/friends_service/internal/storage/neo4j"
)

type ServiceFriends struct {
	log     *slog.Logger
	Storage *neo_storage.Storage
}

func New(log *slog.Logger, storage *neo_storage.Storage) *ServiceFriends {
	return &ServiceFriends{
		log:     log,
		Storage: storage,
	}
}

func (s *ServiceFriends) AddFriend(ctx context.Context, UID, friendID int64) error {
	const op = "friendssrvc.AddFriend"

	log := s.log.With(slog.String("op", op), slog.Int64("UID", UID), slog.Int64("FID", friendID))
	log.Info("Attempting to add friend")
	err := s.Storage.AddFriend(ctx, UID, friendID)
	if err != nil {
		s.log.Warn("ERROR", sl.Err(err))

		return err
	}

	return nil
}

func (s *ServiceFriends) RemoveFriend(ctx context.Context, UID, friendID int64) error {
	const op = "friendssrvc.RemoveFriend"

	log := s.log.With(slog.String("op", op), slog.Int64("UID", UID), slog.Int64("FID", friendID))
	log.Info("Attempting to remove friend")
	err := s.Storage.RemoveFriend(ctx, UID, friendID)
	if err != nil {
		s.log.Warn("ERROR", sl.Err(err))

		return err
	}

	return nil
}

func (s *ServiceFriends) ListUserFriends(ctx context.Context, UID int64) ([]int64, error) {
	const op = "friendssrvc.ListUserFriends"

	log := s.log.With(slog.String("op", op), slog.Int64("UID", UID))
	log.Info("Attempting to list friends")
	FIDs, err := s.Storage.ListFriends(ctx, UID)
	if err != nil {
		s.log.Warn("ERROR", sl.Err(err))

		return nil, err
	}

	return FIDs, nil
}
