package app

import (
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
	grpcapp "github.com/liriquew/social-todo/friends_service/internal/app/app"
	friendssrvc "github.com/liriquew/social-todo/friends_service/internal/grpc/friendsservice"
	"github.com/liriquew/social-todo/friends_service/internal/lib/config"
	neo_storage "github.com/liriquew/social-todo/friends_service/internal/storage/neo4j"
)

type App struct {
	GRPCServer *grpcapp.App
	closers    []func() error
	log        *slog.Logger
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := neo_storage.New(cfg.Neo4jCfg)
	if err != nil {
		panic(err)
	}

	friends_service := friendssrvc.New(log, storage)

	app := grpcapp.New(log, friends_service, cfg.Port)

	mainApp := &App{GRPCServer: app, log: log}
	mainApp.closers = append(mainApp.closers, storage.Close)
	return mainApp
}

func (a *App) Stop() {

	const op = "app.App.Stop"

	for _, c := range a.closers {
		if err := c(); err != nil {
			a.log.Warn("ERROR", sl.Err(fmt.Errorf("%s: %w", op, err)))
		}
	}

	a.GRPCServer.Stop()
}
