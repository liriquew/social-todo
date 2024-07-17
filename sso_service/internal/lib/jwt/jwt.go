package jwt

import (
	"sso_service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var Secret string

// NewToken creates new JWT token for given user and app.
func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
