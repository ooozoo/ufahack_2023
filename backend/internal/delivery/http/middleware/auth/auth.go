package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"ufahack_2023/internal/domain"
	"ufahack_2023/pkg/jwt"
	"ufahack_2023/pkg/logger/sl"
)

const (
	authErrorKey = "auth_error"
	uidKey       = "auth_uid"
	isAdminKey   = "auth_is_admin"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

type PermissionProvider interface {
	IsAdmin(ctx context.Context, uid domain.ID) (bool, error)
}

func New(
	log *slog.Logger,
	secret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "http.middleware.auth"

	log = log.With(
		sl.Op(op),
	)
	log.Info("auth middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt.ParseToken(tokenStr, secret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(r.Context(), authErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UID)
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), authErrorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), uidKey, claims.UID)
			ctx = context.WithValue(ctx, isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
