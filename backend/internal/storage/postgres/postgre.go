package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"ufahack_2023/internal/config"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	"ufahack_2023/internal/storage/models"
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

func (s *Storage) SaveUser(ctx context.Context, username string, passHash []byte) (domain.ID, error) {
	const op = "storage.postgres.Storage.SaveUser"

	stmt, err := s.db.PreparexContext(ctx, "select insert_user(name => $1, password => $2)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	var uid uuid.UUID

	if err := stmt.GetContext(ctx, &uid, username, passHash); err != nil {
		var psqlErr pgx.PgError
		if errors.As(err, &psqlErr) && psqlErr.Code == pgerrcode.UniqueViolation {
			return uuid.Nil, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	const op = "storage.postgres.Storage.GetUserByUsername"

	stmt, err := s.db.PreparexContext(ctx, "select * from get_user_by_username($1)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	if err := stmt.GetContext(ctx, &user, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return models.UserModelToDomain(&user), nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID domain.ID) (bool, error) {
	const op = "storage.postgres.Storage.IsAdmin"

	stmt, err := s.db.PreparexContext(ctx, "select * from is_admin($1)")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var isAdmin sql.NullBool
	if err := stmt.GetContext(ctx, &isAdmin, userID); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if !isAdmin.Valid {
		return false, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return isAdmin.Bool, nil
}
