package authservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/auth"
	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (as *AuthService) NewTokens(user models.User, duration time.Duration) (string, string, error) {
	claims := &auth.Claims{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString(as.jwt.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()

	return accessTokenString, refreshToken, nil
}

func (as *AuthService) ParseToken(ctx context.Context, token string) (int64, time.Time, error) {
	t, err := jwt.ParseWithClaims(token, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.jwt.Secret, nil
	})
	if err != nil {
		return 0, time.Time{}, err
	}

	if !t.Valid {
		return 0, time.Time{}, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*auth.Claims)
	if !ok {
		return 0, time.Time{}, errors.New("invalid claims")
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return 0, time.Time{}, err
	}

	return claims.Id, expiresAt.Time, nil
}
