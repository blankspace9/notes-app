package authservice

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (as *AuthService) NewTokens(user models.User, duration time.Duration) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()

	return accessTokenString, refreshToken, nil
}

func (as *AuthService) ParseToken(ctx context.Context, token string) (int64, time.Time, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.jwtSecret, nil
	})
	if err != nil {
		return 0, time.Time{}, err
	}

	if !t.Valid {
		return 0, time.Time{}, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, time.Time{}, errors.New("invalid claims")
	}

	subject := claims.Subject

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, time.Time{}, errors.New("invalid subject (not a number)")
	}

	return int64(id), claims.ExpiresAt.Time, nil
}
