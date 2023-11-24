package auth

import (
	"context"
	"ufahack_2023/internal/domain"
)

type AuthService interface {
	Login(
		ctx context.Context,
		username string,
		password string,
	) (token string, err error)

	Register(
		ctx context.Context,
		username string,
		password string,
	) (userID domain.ID, err error)

	IsAdmin(
		ctx context.Context,
		userID domain.ID,
	) (bool, error)
}
