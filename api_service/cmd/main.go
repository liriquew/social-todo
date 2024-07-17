package main

import (
	"api_service/internal/app"
	"api_service/internal/lib/config"
	"api_service/internal/lib/logger"
	"api_service/internal/lib/logger/sl"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupPrettySlog()

	log.Info("LOL PON", slog.Any("cfg", cfg))

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
		log.Warn("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}
