package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

type PostgresConnectionInfo struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

// New creates a new instance of the PostgreSQL storage.
func New(connectionInfo PostgresConnectionInfo) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		connectionInfo.Host, connectionInfo.Port, connectionInfo.Username, connectionInfo.DBName, connectionInfo.SSLMode, connectionInfo.Password))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrTokenExists   = errors.New("refresh token already exists")
	ErrTokenNotFound = errors.New("refresh token not fount")
)
