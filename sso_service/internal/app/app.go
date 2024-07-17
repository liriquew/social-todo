package app

import (
	"log/slog"
	"sso_service/internal/app/grpcapp"
	"sso_service/internal/lib/config"
	"sso_service/internal/sevices/auth"
	"sso_service/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	auth := auth.New(log, storage)

	app := grpcapp.New(log, auth, cfg.Port)
	return &App{GRPCServer: app}
}
