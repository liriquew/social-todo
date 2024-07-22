package auth

import (
	"errors"
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
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	token, err := a.authClient.Login(c, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, auth_grpc.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		if errors.Is(err, auth_grpc.ErrInvalidArgument) {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusInternalServerError)
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
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	uid, err := a.authClient.Register(c, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, auth_grpc.ErrAlreadyExists) {
			c.Status(http.StatusConflict)
			return
		}
		if errors.Is(err, auth_grpc.ErrInvalidArgument) {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("uid: %d", uid))
}
