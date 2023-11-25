package model

import "ufahack_2023/internal/domain"

type Subject struct {
	ID   domain.ID `db:"subjects_uid"`
	Name string    `db:"subject_name"`
}

func SubjectModelToDomain(s *Subject) *domain.Subject {
	return &domain.Subject{
		ID:   s.ID,
		Name: s.Name,
	}
}

func SubjectsListToDomain(sl []*Subject) []*domain.Subject {
	var sdl []*domain.Subject
	for _, s := range sl {
		sdl = append(sdl, SubjectModelToDomain(s))
	}

	return sdl
}
