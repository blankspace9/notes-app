package models

import "time"

type Note struct {
	ID        int64     `json:"id"`
	Note      string    `json:"note"`
	UserID    int64     `json:"userID,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type NoteRequest struct {
	Note string `json:"note"`
}
