package login

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"ufahack_2023/internal/delivery/http/handlers/common"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/service/auth"
	resp "ufahack_2023/pkg/api/response"
	"ufahack_2023/pkg/logger/sl"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserLoginer
type UserLoginer interface {
	Login(ctx context.Context, username string, password string) (*domain.User, string, error)
}

func New(log *slog.Logger, loginer UserLoginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.auth.login"

		log := log.With(
			sl.Op(op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		common.DecodeRequest(log, w, r, &req)
		common.ValidateRequest(log, w, r, req)

		user, token, err := loginer.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				log.Warn("invalid username or password")

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid username or password"))

				return
			}

			log.Error("failed to login user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to login"))

			return
		}

		log.Debug("user successfully logged in", slog.String("id", user.ID.String()))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Token:    token,
		})
	}
}
