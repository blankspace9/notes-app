package authservice

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/blankspace9/notes-app/internal/config"
	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
	"github.com/blankspace9/notes-app/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type AuthService struct {
	log          *slog.Logger
	userManager  UserManager
	tokenManager TokenManager
	jwt          config.JWT
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUserExists         = errors.New("user already exists")

	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type UserManager interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserById(ctx context.Context, id int64) (models.User, error)
}

type TokenManager interface {
	SaveToken(ctx context.Context, token models.Token) error
	GetToken(ctx context.Context, token string) (models.Token, error)
	UpdateToken(ctx context.Context, token models.Token) error
}

// New return a new instance of the Auth service
func New(log *slog.Logger, userManager UserManager, tokenManager TokenManager, tokens config.JWT) *AuthService {
	return &AuthService{
		log:          log,
		userManager:  userManager,
		tokenManager: tokenManager,
		jwt:          tokens,
	}
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
		if errors.Is(err, storage.ErrUserExists) {
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
	const op = "services.AuthService.Login"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to login user")

	user, err := as.userManager.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
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

	accessToken, refreshToken, err := as.NewTokens(user, as.jwt.AccessTokenTTL)
	if err != nil {
		as.log.Error("failed to generate tokens", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	token := models.Token{
		UserId:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(as.jwt.RefreshTokenTTL),
	}
	err = as.tokenManager.SaveToken(ctx, token)
	if err != nil {
		if errors.Is(err, storage.ErrTokenExists) {
			err = as.tokenManager.UpdateToken(ctx, token)
			if err != nil {
				as.log.Error("failed to update token", sl.Err(err))

				return "", "", err
			}
		} else {
			as.log.Error("failed to save token", sl.Err(err))

			return "", "", err
		}
	}

	log.Info("user logged in successfully")

	return accessToken, refreshToken, nil
}

func (as *AuthService) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	const op = "services.AuthService.RefreshTokens"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to login user")

	session, err := as.tokenManager.GetToken(ctx, token)
	if err != nil {
		if errors.Is(err, storage.ErrTokenNotFound) {
			as.log.Warn("refresh token not found", sl.Err(err))

			return "", "", fmt.Errorf("%s: %w", op, ErrRefreshTokenNotFound)
		}

		as.log.Error("failed to get token", sl.Err(err))

		return "", "", err
	}

	user, err := as.userManager.UserById(ctx, session.UserId)
	if err != nil {
		as.log.Error("failed to get user", sl.Err(err))

		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		as.log.Error("failed to refresh tokens", sl.Err(ErrRefreshTokenExpired))

		return "", "", ErrRefreshTokenExpired
	}

	accessToken, refreshToken, err := as.NewTokens(user, as.jwt.AccessTokenTTL)
	if err != nil {
		as.log.Error("failed to generate tokens", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	t := models.Token{
		UserId:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(as.jwt.RefreshTokenTTL),
	}
	err = as.tokenManager.SaveToken(ctx, t)
	if err != nil {
		if errors.Is(err, storage.ErrTokenExists) {
			err = as.tokenManager.UpdateToken(ctx, t)
			if err != nil {
				as.log.Error("failed to update token", sl.Err(err))

				return "", "", err
			}
		} else {
			as.log.Error("failed to save token", sl.Err(err))

			return "", "", err
		}
	}

	log.Info("tokens refreshed successfully")

	return accessToken, refreshToken, nil
}
