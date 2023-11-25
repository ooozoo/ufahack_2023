package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"ufahack_2023/internal/domain"
	resp "ufahack_2023/pkg/api/response"
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

func NewAuth(
	log *slog.Logger,
	secret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "http.middleware.auth.Auth"

	log = log.With(
		sl.Op(op),
	)
	log.Info("auth middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log = log.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

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
				ctx = context.WithValue(ctx, uidKey, claims.UID)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), uidKey, claims.UID)
			ctx = context.WithValue(ctx, isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func NewAdminOnly(log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "http.middleware.auth.AdminOnly"

	log = log.With(
		sl.Op(op),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log = log.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			isAdmin, ok := r.Context().Value(isAdminKey).(bool)
			if !ok || !isAdmin {
				uid, ok := r.Context().Value(uidKey).(uuid.UUID)
				if !ok {
					log.Warn("anonymous tried to access admin endpoint")
				} else {
					log.Warn("user tried to access admin endpoint", slog.String("uid", uid.String()))
				}
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("action allowed only for admins"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NewUserOnly(log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "http.middleware.auth.UserOnly"

	log = log.With(
		sl.Op(op),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log = log.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			uid, ok := r.Context().Value(uidKey).(uuid.UUID)
			if !ok || uid == uuid.Nil {
				log.Warn("anonymous tried to access user endpoints", slog.String("request_id", middleware.GetReqID(r.Context())))
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, resp.Error("action allowed only for users"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
