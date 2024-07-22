package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_grpc "github.com/liriquew/social-todo/api_service/internal/clients/authgrpc"
)

func (a *Auth) AuthRequired(c *gin.Context) {
	a.log.Info("middleware")

	token := c.GetHeader("Authorization")
	if token == "" {
		c.String(http.StatusUnauthorized, "jwt token required")
	}

	uid, err := a.authClient.Authorize(c, token)
	if err != nil {
		if errors.Is(err, auth_grpc.ErrMissJWTToken) {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.Set("uid", uid)

	c.Next()
}
