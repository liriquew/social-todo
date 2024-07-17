package app

import (
	apiapp "api_service/internal/app/app"
	auth_grpc "api_service/internal/clients/authgrpc"
	"api_service/internal/lib/config"
	"api_service/internal/rest/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type App struct {
	GinRouter *gin.Engine
}

func New(log *slog.Logger, cfg config.Config) *App {
	authClient, err := auth_grpc.New(log, cfg.ClientGRPC.AuthPort, cfg.ClientGRPC.Timeout, cfg.ClientGRPC.Retries)
	if err != nil {
		panic(err)
	}

	auth := auth.New(log, authClient)
	r := apiapp.New(auth)

	return &App{r}
}
