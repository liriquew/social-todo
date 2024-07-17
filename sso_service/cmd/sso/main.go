package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso_service/internal/app"
	"sso_service/internal/lib/config"
	"sso_service/internal/lib/jwt"
	"sso_service/internal/lib/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupPrettySlog()

	log.Info("", slog.Any("CONFIG", cfg))

	jwt.Secret = cfg.JWTSecret

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
