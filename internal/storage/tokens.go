package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/lib/pq"
)

func (s *Storage) SaveToken(ctx context.Context, token models.Token) error {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO tokens (user_id, token, expires_at) values ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, token.UserId, token.Token, token.ExpiresAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrTokenExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetToken(ctx context.Context, token string) (models.Token, error) {
	const op = "storage.postgres.GetToken"

	stmt, err := s.db.Prepare("SELECT id, user_id, token, expires_at FROM tokens WHERE token=$1")
	if err != nil {
		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, token)

	var t models.Token
	err = row.Scan(&t.ID, &t.UserId, &t.Token, &t.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Token{}, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}

		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (s *Storage) UpdateToken(ctx context.Context, token models.Token) error {
	const op = "storage.postgres.UpdateToken"

	stmt, err := s.db.Prepare("UPDATE tokens SET token=$1, expires_at=$2 WHERE user_id=$3")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, token.Token, token.ExpiresAt, token.UserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
