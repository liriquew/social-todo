package apiapp

import (
	"github.com/liriquew/social-todo/api_service/internal/rest/auth"
	"github.com/liriquew/social-todo/api_service/internal/rest/friends"
	"github.com/liriquew/social-todo/api_service/internal/rest/notes"
	handlers "github.com/liriquew/social-todo/api_service/internal/rest/other"

	"github.com/gin-gonic/gin"
)

func New(auth auth.AuthAPI, notes notes.NotesAPI, friends friends.FriendsAPI, other handlers.GeneralAPI) *gin.Engine {
	r := gin.New()

	r.POST("/signin", auth.Login)
	r.POST("/signup", auth.Register)

	notesAPI := r.Group("/note")
	notesAPI.Use(auth.AuthRequired)
	{
		notesAPI.GET("/listid", notes.ListIDs)
		notesAPI.GET("/listnotes", notes.ListNotes)
		notesAPI.POST("/", notes.Create)
		notesAPI.GET("/:id", notes.Get)
		notesAPI.PATCH("/:id", notes.Update)
		notesAPI.DELETE("/:id", notes.Delete)
	}

	friendsAPI := r.Group("/friends")
	friendsAPI.Use(auth.AuthRequired)
	{
		friendsAPI.POST("/add", friends.AddFriend)
		friendsAPI.POST("/remove", friends.RemoveFriend)
		friendsAPI.GET("/list", friends.ListFriend)
	}

	news := r.Group("/news")
	news.Use(auth.AuthRequired)
	{
		news.GET("", other.ListLastNotes)
	}

	return r
}
