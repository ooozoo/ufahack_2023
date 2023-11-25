package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"ufahack_2023/internal/domain"
)

type User struct {
	UID       uuid.UUID    `db:"users_uid"`
	Name      string       `db:"user_name"`
	Pass      []byte       `db:"pass"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	IsAdmin   bool         `db:"is_admin"`
	IsDeleted bool         `db:"is_deleted"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func UserModelToDomain(user *User) *domain.User {
	return &domain.User{
		ID:           user.UID,
		Username:     user.Name,
		PasswordHash: user.Pass,
	}
}
