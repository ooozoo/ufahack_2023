package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"ufahack_2023/internal/config"
)

type Storage struct {
	db *sqlx.DB
}

func New(databaseConfig config.DatabaseConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.User,
		databaseConfig.Password,
		databaseConfig.Name,
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Storage.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
