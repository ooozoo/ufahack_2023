package register

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"ufahack_2023/internal/delivery/http/handlers/common"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	resp "ufahack_2023/pkg/api/response"
	"ufahack_2023/pkg/logger/sl"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	UserID domain.ID `json:"user_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserRegister
type UserRegister interface {
	Register(ctx context.Context, username string, password string) (domain.ID, error)
}

func New(log *slog.Logger, register UserRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.auth.register"

		log := log.With(
			sl.Op(op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if !common.DecodeRequest(log, w, r, &req) {
			return
		}
		if !common.ValidateRequest(log, w, r, req) {
			return
		}

		uid, err := register.Register(r.Context(), req.Username, req.Password)
		if err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				log.Warn("user already exists", slog.String("username", req.Username))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("user already exists"))

				return
			}

			log.Error("failed to register user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to register user"))

			return
		}

		log.Info("registered user", slog.String("id", uid.String()))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			UserID:   uid,
		})
	}
}
