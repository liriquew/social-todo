package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/social-todo/api_service/pkg/logger"
	"github.com/liriquew/social-todo/api_service/pkg/logger/sl"

	"github.com/liriquew/social-todo/api_service/internal/app"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupPrettySlog()

	log.Info("CONFIG", slog.Any("cfg", cfg))

	r := app.New(log, *cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r.GinRouter,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("listen", sl.Err(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Warn("Server forced to shutdown:", sl.Err(err))
	}

	log.Info("Server exiting")
}
