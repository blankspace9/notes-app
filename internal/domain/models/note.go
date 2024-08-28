package models

import "time"

type Note struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	UserID    int64
	CreatedAt time.Time `json:"createdAt"`
}
