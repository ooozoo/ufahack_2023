package subject

import (
	"ufahack_2023/internal/domain"
)

type Subject struct {
	ID   domain.ID `json:"id"`
	Name string    `json:"name"`
}

func ToModelOne(s *domain.Subject) *Subject {
	return &Subject{
		ID:   s.ID,
		Name: s.Name,
	}
}

func ToModelMany(sl []*domain.Subject) []*Subject {
	var subjects []*Subject
	for _, s := range sl {
		subjects = append(subjects, ToModelOne(s))
	}

	return subjects
}
