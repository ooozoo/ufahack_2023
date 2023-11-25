package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"

	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	"ufahack_2023/internal/storage/model"
)

func (s *Storage) SaveUser(ctx context.Context, username string, passHash []byte) (domain.ID, error) {
	const op = "storage.postgres.Storage.SaveUser"

	stmt, err := s.db.PreparexContext(ctx, "select insert_user($1, $2)")
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

	var user model.User
	if err := stmt.GetContext(ctx, &user, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model.UserModelToDomain(&user), nil
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
