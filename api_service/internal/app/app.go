package app

import (
	"log/slog"

	apiapp "github.com/liriquew/social-todo/api_service/internal/app/app"
	auth_grpc "github.com/liriquew/social-todo/api_service/internal/clients/authgrpc"
	notes_grpc "github.com/liriquew/social-todo/api_service/internal/clients/notesgrpc"
	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/social-todo/api_service/internal/rest/auth"
	"github.com/liriquew/social-todo/api_service/internal/rest/notes"

	"github.com/gin-gonic/gin"
)

type App struct {
	GinRouter *gin.Engine
}

func New(log *slog.Logger, cfg config.Config) *App {
	authClient, err := auth_grpc.New(log, cfg.AuthConfig)
	if err != nil {
		panic(err)
	}
	notesClient, err := notes_grpc.New(log, cfg.NoteConfig)

	auth := auth.New(log, authClient)
	notes := notes.New(log, notesClient)

	r := apiapp.New(auth, notes)

	return &App{r}
}
