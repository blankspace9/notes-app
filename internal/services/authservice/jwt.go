package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/auth"
	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (as *AuthService) NewTokens(user models.User, duration time.Duration) (string, string, error) {
	const op = "services.AuthService.NewTokens"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to generate tokens for user", slog.Int64("userID", user.ID))

	claims := &auth.Claims{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString(as.jwt.Secret)
	if err != nil {
		log.Error("failed to signed token", sl.Err(err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken := uuid.New().String()

	log.Debug("tokens generated successfully", slog.String("accessToken", accessTokenString))

	return accessTokenString, refreshToken, nil
}

func (as *AuthService) ParseToken(ctx context.Context, token string) (int64, time.Time, error) {
	const op = "services.AuthService.ParseToken"

	log := as.log.With(slog.String("op", op))

	log.Info("attempting to parse token")

	t, err := jwt.ParseWithClaims(token, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("failed to signed token")

			return nil, fmt.Errorf("%s: %s %s", op, "unexpected signing method: ", token.Header["alg"])
		}
		return as.jwt.Secret, nil
	})
	if err != nil {
		log.Error("failed to parse token", sl.Err(err))

		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	if !t.Valid {
		log.Error("invalid token")

		return 0, time.Time{}, fmt.Errorf("%s: %w", op, errors.New("invalid token"))
	}

	claims, ok := t.Claims.(*auth.Claims)
	if !ok {
		log.Error("invalid claims")

		return 0, time.Time{}, fmt.Errorf("%s: %w", op, errors.New("invalid claims"))
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		log.Error("invalid claims", sl.Err(err))

		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return claims.Id, expiresAt.Time, nil
}
