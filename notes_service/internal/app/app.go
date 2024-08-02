package app

import (
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"
	grpcapp "github.com/liriquew/social-todo/notes_service/internal/app/app"
	"github.com/liriquew/social-todo/notes_service/internal/grpc/notessrvc"
	"github.com/liriquew/social-todo/notes_service/internal/lib/config"
	"github.com/liriquew/social-todo/notes_service/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
	closers    []func() error
	log        *slog.Logger
}

func New(log *slog.Logger, cfg config.Config) *App {
	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	notes_service := notessrvc.New(log, storage)

	app := grpcapp.New(log, notes_service, cfg.Port)

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
