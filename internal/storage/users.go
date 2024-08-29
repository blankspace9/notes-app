package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/lib/pq"
)

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash, created_at) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email, passHash, time.Now())

	var uid int64
	err = row.Scan(&uid)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.UserByEmail"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash, created_at FROM users WHERE email=$1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserById(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.postgres.UserById"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash, created_at FROM users WHERE id=$1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
