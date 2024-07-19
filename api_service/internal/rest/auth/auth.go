package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/liriquew/social-todo/api_service/internal/models"

	auth_grpc "github.com/liriquew/social-todo/api_service/internal/clients/authgrpc"

	"github.com/gin-gonic/gin"
)

type AuthAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	AuthRequired(c *gin.Context)
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

func (a *Auth) AuthRequired(c *gin.Context) {
	a.log.Info("middleware")

	token := c.GetHeader("Authorization")
	if token == "" {
		c.String(http.StatusUnauthorized, "jwt token required")
	}

	a.authClient.Authorize(c, token)

	// TODO: обратиться в сервис заметок

	// TODO: обработать то, что вернул сервис заметок
}
