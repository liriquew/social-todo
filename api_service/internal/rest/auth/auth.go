package auth

import (
	"api_service/internal/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type Auth struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Auth {
	return &Auth{
		log: log,
	}
}

func (a *Auth) Login(c *gin.Context) {
	a.log.Info("Login")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

func (a *Auth) Register(c *gin.Context) {
	a.log.Info("Register")
}
