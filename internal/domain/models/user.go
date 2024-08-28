package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	PassHash  []byte    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
