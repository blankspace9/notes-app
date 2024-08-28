package authservice

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/blankspace9/notes-app/internal/config"
	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/blankspace9/notes-app/internal/lib/jwt"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type AuthService struct {
	log          *slog.Logger
	userManager  UserManager
	tokenManager TokenManager
	tokens       config.Tokens
	jwtSecret    string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUserExists         = errors.New("user already exists")
)

type UserManager interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
}

type TokenManager interface {
	SaveToken(ctx context.Context, token models.Token) error
}

func (as *AuthService) RegisterUser(ctx context.Context, email, password string) (int64, error) {
	const op = "services.AuthService.RegisterUser"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := as.userManager.SaveUser(ctx, email, passHash)
	if err != nil {
		// TODO: change errors.New to storage.Error
		if errors.Is(err, errors.New("storageError (user already exists)")) {
			log.Warn("user already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	const op = "auth.Login"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to login user")

	user, err := as.userManager.UserByEmail(ctx, email)
	// TODO: change errors.New to storage.Error
	if err != nil {
		if errors.Is(err, errors.New("storageError (user not found)")) {
			as.log.Warn("user not found", sl.Err(err))

			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		as.log.Error("failed to get user", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		as.log.Warn("invalid credentials", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	accessToken, refreshToken, err := jwt.NewTokens(user, as.tokens.AccessTokenTTL, as.jwtSecret)
	if err != nil {
		as.log.Error("failed to generate tokens", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := as.tokenManager.SaveToken(ctx, models.Token{
		UserId:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(as.tokens.RefreshTokenTTL),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
