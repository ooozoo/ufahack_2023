package subject

import (
	"context"
	"log/slog"

	"ufahack_2023/internal/domain"
)

type Provider interface {
	ListSubjects(ctx context.Context) ([]*domain.Subject, error)
	GetSubjectByID(ctx context.Context, id domain.ID) (*domain.Subject, error)
}

type Saver interface {
	SaveSubject(ctx context.Context, name string) (domain.ID, error)
	UpdateSubject(ctx context.Context, id domain.ID, name string) (*domain.Subject, error)
	DeleteSubject(ctx context.Context, id domain.ID) error
}

type Subject struct {
	log      *slog.Logger
	provider Provider
	saver    Saver
}

func New(
	log *slog.Logger,
	provider Provider,
	saver Saver) *Subject {
	return &Subject{
		log:      log,
		provider: provider,
		saver:    saver,
	}
}

func (s *Subject) ListSubjects(ctx context.Context) ([]*domain.Subject, error) {
	return s.provider.ListSubjects(ctx)
}

func (s *Subject) GetSubjectByID(ctx context.Context, id domain.ID) (*domain.Subject, error) {
	return s.provider.GetSubjectByID(ctx, id)
}

func (s *Subject) SaveSubject(ctx context.Context, name string) (domain.ID, error) {
	return s.saver.SaveSubject(ctx, name)
}

func (s *Subject) UpdateSubject(ctx context.Context, id domain.ID, name string) (*domain.Subject, error) {
	return s.saver.UpdateSubject(ctx, id, name)
}

func (s *Subject) DeleteSubject(ctx context.Context, id domain.ID) error {
	return s.saver.DeleteSubject(ctx, id)
}
