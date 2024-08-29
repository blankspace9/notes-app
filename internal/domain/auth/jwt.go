package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Id int64
	jwt.RegisteredClaims
}
