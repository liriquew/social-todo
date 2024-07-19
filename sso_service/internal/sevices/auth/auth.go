package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/liriquew/social-todo/sso_service/internal/lib/jwt"
	"github.com/liriquew/social-todo/sso_service/internal/models"
	"github.com/liriquew/social-todo/sso_service/internal/storage"

	"github.com/liriquew/social-todo/api_service/internal/lib/logger/sl"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log     *slog.Logger
	Storage StorageProvider
}

type StorageProvider interface {
	SaveUser(context.Context, string, []byte) (int64, error)
	User(context.Context, string) (models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist          = errors.New("user already exist")
	ErrUserNotFound       = errors.New("user not found")
)

func New(log *slog.Logger, Storage StorageProvider) *Auth {
	return &Auth{
		log:     log,
		Storage: Storage,
	}
}

func (a *Auth) Login(ctx context.Context, username, password string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("username", username))
	log.Info("attempting to login user")

	user, err := a.Storage.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, username, password string) (int64, error) {
	const op = "auth.Register"

	log := a.log.With(slog.String("op", op), slog.String("username", username))
	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Warn("failed to generate hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := a.Storage.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			a.log.Warn("user already exist")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}
		a.log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, err
}
