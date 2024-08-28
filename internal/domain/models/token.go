package models

import (
	"time"
)

type Token struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
