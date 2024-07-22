package apiapp

import (
	"github.com/liriquew/social-todo/api_service/internal/rest/auth"
	"github.com/liriquew/social-todo/api_service/internal/rest/notes"

	"github.com/gin-gonic/gin"
)

func New(auth auth.AuthAPI, notes notes.NotesAPI) *gin.Engine {
	r := gin.New()

	r.POST("/signin", auth.Login)
	r.POST("/signup", auth.Register)

	notesAPI := r.Group("/note")
	notesAPI.Use(auth.AuthRequired)
	{
		notesAPI.POST("/", notes.Create)
		notesAPI.GET("/:id", notes.Get)
		notesAPI.PATCH("/:id", notes.Update)
		notesAPI.DELETE("/:id", notes.Delete)
	}

	return r
}
