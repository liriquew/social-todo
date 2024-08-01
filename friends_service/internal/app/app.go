package app

import (
	"log/slog"

	grpcapp "github.com/liriquew/social-todo/friends_service/internal/app/app"
	friendssrvc "github.com/liriquew/social-todo/friends_service/internal/grpc/friendsservice"
	"github.com/liriquew/social-todo/friends_service/internal/lib/config"
	neo_storage "github.com/liriquew/social-todo/friends_service/internal/storage/neo4j"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := neo_storage.New(cfg.Neo4jCfg)
	if err != nil {
		panic(err)
	}

	friends_service := friendssrvc.New(log, storage)

	app := grpcapp.New(log, friends_service, cfg.Port)
	return &App{GRPCServer: app}
}
