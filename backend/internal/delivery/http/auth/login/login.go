package login

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"ufahack_2023/internal/domain"
	resp "ufahack_2023/internal/lib/api/response"
	"ufahack_2023/internal/lib/logger/sl"
	"ufahack_2023/internal/service/auth"
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
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "delivery.http.auth.login.New"

		log := log.With(
			sl.Op(op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("empty request"))

				return
			}

			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Debug("request body decoded")

		if err := v.Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		user, token, err := loginer.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				log.Error("invalid username or password")

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.Error("invalid username or password"))

				return
			}

			log.Error("failed to login user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to login"))

			return
		}

		log.Info("user successfully logged in", slog.String("id", user.ID.String()))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Token:    token,
		})
	}
}
