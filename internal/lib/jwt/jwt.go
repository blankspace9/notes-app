package jwt

import (
	"time"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func NewTokens(user models.User, duration time.Duration, jwtSecret string) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()

	return accessTokenString, refreshToken, nil
}
