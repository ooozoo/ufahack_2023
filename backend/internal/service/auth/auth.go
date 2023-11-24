package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	"ufahack_2023/pkg/jwt"
	"ufahack_2023/pkg/logger/sl"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		passHash []byte,
	) (uid domain.ID, err error)
}

type UserProvider interface {
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*domain.User, error)

	IsAdmin(
		ctx context.Context,
		userID domain.ID,
	) (bool, error)
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	secret       string
	tokenTTL     time.Duration
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	secret string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		secret:       secret,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	username string,
	pass string,
) (domain.ID, error) {
	const op = "service.auth.Auth.Register"

	log := a.log.With(
		sl.Op(op),
		slog.String("username", username),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return uuid.Nil, err
	}

	id, err := a.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return uuid.Nil, err
	}

	return id, nil
}

func (a *Auth) Login(
	ctx context.Context,
	username string,
	pass string,
) (*domain.User, string, error) {
	const op = "service.auth.Auth.Login"

	log := a.log.With(
		sl.Op(op),
		slog.String("username", username),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return nil, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass)); err != nil {
		a.log.Warn("invalid credentials", sl.Err(err))

		return nil, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user.ID, a.secret, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return user, token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID domain.ID) (bool, error) {
	const op = "service.auth.Auth.IsAdmin"

	log := a.log.With(
		sl.Op(op),
		slog.String("user_id", userID.String()),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed to check if user is admin", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
