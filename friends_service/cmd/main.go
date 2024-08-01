package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/liriquew/social-todo/friends_service/internal/app"
	"github.com/liriquew/social-todo/friends_service/internal/lib/config"

	"github.com/liriquew/social-todo/api_service/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupPrettySlog()

	log.Info("", slog.Any("CONFIG", cfg))

	application := app.New(log, cfg)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.GRPCServer.Stop()

	log.Info("Gracefully stopped")
}
