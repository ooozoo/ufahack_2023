package model

import (
	"github.com/google/uuid"

	"ufahack_2023/internal/domain"
)

type User struct {
	UID     uuid.UUID `db:"users_uid"`
	Name    string    `db:"user_name"`
	Pass    []byte    `db:"pass"`
	IsAdmin bool      `db:"is_admin"`
}

func UserModelToDomain(user *User) *domain.User {
	return &domain.User{
		ID:           user.UID,
		Username:     user.Name,
		PasswordHash: user.Pass,
	}
}
