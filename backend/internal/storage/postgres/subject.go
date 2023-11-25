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

func (s *Storage) ListSubjects(ctx context.Context) ([]*domain.Subject, error) {
	const op = "storage.postgres.ListSubjects"

	stmt, err := s.db.PreparexContext(ctx, "select get_subjects()")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var subjects []*model.Subject
	if err := stmt.SelectContext(ctx, subjects); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model.SubjectsListToDomain(subjects), err
}

func (s *Storage) GetSubjectByID(ctx context.Context, id domain.ID) (*domain.Subject, error) {
	const op = "storage.postgres.GetSubjects"

	stmt, err := s.db.PreparexContext(ctx, "select get_subject_by_id($1)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var subject *model.Subject
	if err := stmt.GetContext(ctx, subject); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model.SubjectModelToDomain(subject), nil
}

func (s *Storage) SaveSubject(ctx context.Context, name string) (domain.ID, error) {
	const op = "storage.postgres.SaveSubject"

	stmt, err := s.db.PreparexContext(ctx, "select insert_subject($1)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	var uid domain.ID
	if err := stmt.GetContext(ctx, &uid, name); err != nil {
		var psqlErr pgx.PgError
		if errors.As(err, &psqlErr) && psqlErr.Code == pgerrcode.UniqueViolation {
			return uuid.Nil, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) UpdateSubject(ctx context.Context, id domain.ID, name string) (*domain.Subject, error) {
	const op = "storage.postgres.UpdateSubject"
	panic("implement me")
}

func (s *Storage) DeleteSubject(ctx context.Context, id domain.ID) error {
	//TODO implement me
	panic("implement me")
}
