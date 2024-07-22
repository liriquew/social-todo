package app

import (
	"log/slog"

	grpcapp "github.com/liriquew/social-todo/notes_service/internal/app/app"
	"github.com/liriquew/social-todo/notes_service/internal/grpc/notessrvc"
	"github.com/liriquew/social-todo/notes_service/internal/lib/config"
	"github.com/liriquew/social-todo/notes_service/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	notes_service := notessrvc.New(log, storage)

	app := grpcapp.New(log, notes_service, cfg.Port)
	return &App{GRPCServer: app}
}
