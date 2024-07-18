package apiapp

import (
	"api_service/internal/rest/auth"

	"github.com/gin-gonic/gin"
)

func New(auth auth.AuthAPI) *gin.Engine {
	r := gin.New()

	r.POST("/signin", auth.Login)
	r.POST("/signup", auth.Register)

	r.GET("/note", auth.AuthRequired)

	return r
}
