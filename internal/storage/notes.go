package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/models"
)

func (s *Storage) SaveNote(ctx context.Context, note models.Note) (int64, error) {
	const op = "storage.postgres.SaveNote"

	stmt, err := s.db.Prepare("INSERT INTO notes(note, user_id, created_at) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, note.Note, note.UserID, time.Now())

	var insertedID int64
	err = row.Scan(&insertedID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return insertedID, nil
}

func (s *Storage) GetNotesByUserId(ctx context.Context, userID int64) ([]models.Note, error) {
	const op = "storage.postgres.GetNotesByUserId"

	stmt, err := s.db.Prepare("SELECT id, note, created_at FROM notes WHERE user_id=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note

		err = rows.Scan(&note.ID, &note.Note, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Storage) GetNotesPageByUserId(ctx context.Context, userID int64, page, limit int) ([]models.Note, error) {
	const op = "storage.postgres.GetNotesByUserId"

	offset := (page - 1) * limit

	stmt, err := s.db.Prepare("SELECT id, note, created_at FROM notes WHERE user_id=$1 LIMIT $2 OFFSET $3")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note

		err = rows.Scan(&note.ID, &note.Note, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}
