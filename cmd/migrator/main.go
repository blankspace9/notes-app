package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/blankspace9/notes-app/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath, migrationsTable string

	cfg := config.MigrateMustLoad()

	postgresPath := fmt.Sprintf("%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migration")
	down := flag.Bool("down", false, "up migrations")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("postgres://%s?sslmode=%s&x-migrations-table=%s",
		postgresPath, cfg.SSLMode, migrationsTable))
	if err != nil {
		panic(err)
	}
	if !(*down) {
		fmt.Println("UP")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}

			panic(err)
		}
	} else {
		fmt.Println("DOWN")
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}

			panic(err)
		}
	}

	fmt.Println("migrations applied successfully")
}
