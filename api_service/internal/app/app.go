package app

import (
	apiapp "api_service/internal/app/app"
	"api_service/internal/rest/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type App struct {
	GinRouter *gin.Engine
}

func New(log *slog.Logger) *App {
	auth := auth.New(log)
	r := apiapp.New(auth)

	return &App{r}
}
