package jwt

import (
	"fmt"
	"sso_service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var Secret string

// NewToken creates new JWT token for given user and app.
func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["uid"] = user.UID

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate проверяет jwt токен и возвращает uid пользователя
// по которому уже можно однозначно определять полльзователя
func Validate(tokenString string) (int64, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})

	fmt.Println(err)
	fmt.Println("valid token:", token.Valid)

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return int64((*claims)["uid"].(float64)), nil // да
}
