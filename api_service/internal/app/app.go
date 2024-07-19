package app

import (
	"log/slog"

	apiapp "github.com/liriquew/social-todo/api_service/internal/app/app"
	auth_grpc "github.com/liriquew/social-todo/api_service/internal/clients/authgrpc"
	"github.com/liriquew/social-todo/api_service/internal/rest/auth"
	"github.com/liriquew/social-todo/api_service/pkg/config"

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
