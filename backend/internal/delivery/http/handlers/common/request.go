package common

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "ufahack_2023/pkg/api/response"
	"ufahack_2023/pkg/api/valid"
	"ufahack_2023/pkg/logger/sl"
)

func DecodeRequest(log *slog.Logger, w http.ResponseWriter, r *http.Request, req any) {
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Warn("request body is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		log.Error("failed to decode request body", sl.Err(err))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}
}

func ValidateRequest(log *slog.Logger, w http.ResponseWriter, r *http.Request, req any) {
	v := valid.GetValidator()

	if err := v.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		log.Error("invalid request", sl.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.ValidationError(validateErr))

		return
	}
}
