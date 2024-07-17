package auth

import (
	"api_service/internal/models"
	"fmt"
	"log/slog"
	"net/http"

	auth_grpc "api_service/internal/clients/authgrpc"

	"github.com/gin-gonic/gin"
)

type AuthAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type Auth struct {
	log        *slog.Logger
	authClient *auth_grpc.Client
}

func New(log *slog.Logger, authClient *auth_grpc.Client) *Auth {
	return &Auth{
		log:        log,
		authClient: authClient,
	}
}

func (a *Auth) Login(c *gin.Context) {
	a.log.Info("Login")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authClient.Login(c, user.Username, user.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a *Auth) Register(c *gin.Context) {
	a.log.Info("Register")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, err := a.authClient.Register(c, user.Username, user.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

	if uid != 0 {
		c.String(http.StatusOK, "register - ok")
	} else {
		c.String(http.StatusOK, fmt.Sprintf("error: %s", err))
	}
}
