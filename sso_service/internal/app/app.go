package app

import (
	"log/slog"

	"github.com/liriquew/social-todo/sso_service/internal/app/grpcapp"
	"github.com/liriquew/social-todo/sso_service/internal/lib/config"
	"github.com/liriquew/social-todo/sso_service/internal/sevices/auth"
	"github.com/liriquew/social-todo/sso_service/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	auth := auth.New(log, storage)

	app := grpcapp.New(log, auth, cfg.Port)
	return &App{GRPCServer: app}
}
