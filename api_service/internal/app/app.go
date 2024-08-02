package app

import (
	"log/slog"

	apiapp "github.com/liriquew/social-todo/api_service/internal/app/app"
	auth_grpc "github.com/liriquew/social-todo/api_service/internal/clients/authgrpc"
	friends_grpc "github.com/liriquew/social-todo/api_service/internal/clients/friendsgrpc"
	notes_grpc "github.com/liriquew/social-todo/api_service/internal/clients/notesgrpc"
	"github.com/liriquew/social-todo/api_service/internal/lib/config"
	"github.com/liriquew/social-todo/api_service/internal/rest/auth"
	"github.com/liriquew/social-todo/api_service/internal/rest/friends"
	"github.com/liriquew/social-todo/api_service/internal/rest/notes"
	handlers "github.com/liriquew/social-todo/api_service/internal/rest/other"

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
	if err != nil {
		panic(err)
	}
	friendsClient, err := friends_grpc.New(log, cfg.FriendsConfig)
	if err != nil {
		panic(err)
	}

	auth := auth.New(log, authClient)
	notes := notes.New(log, notesClient)
	friends := friends.New(log, friendsClient)
	other := handlers.New(log, notesClient, friendsClient)

	r := apiapp.New(auth, notes, friends, other)

	return &App{r}
}
