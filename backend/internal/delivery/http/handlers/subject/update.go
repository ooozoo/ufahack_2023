package subject

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"ufahack_2023/internal/delivery/http/handlers/common"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	resp "ufahack_2023/pkg/api/response"
	"ufahack_2023/pkg/logger/sl"
)

type Updater interface {
	UpdateSubject(ctx context.Context, id domain.ID, Name string) (*domain.Subject, error)
}

type UpdateRequest struct {
	ID   domain.ID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

type UpdateResponse struct {
	resp.Response
	Subject *Subject
}

func NewUpdateSubject(
	log *slog.Logger,
	updater Updater,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.subject.UpdateSubject"

		log = log.With(
			sl.Op(op),
		)

		var req UpdateRequest

		common.DecodeRequest(log, w, r, &req)
		common.ValidateRequest(log, w, r, req)

		subject, err := updater.UpdateSubject(r.Context(), req.ID, req.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Warn("subject not exists", slog.String("uid", req.ID.String()))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("subject not exists"))

				return
			}

			if errors.Is(err, storage.ErrAlreadyExists) {
				log.Warn("subject name collision", slog.String("uid", req.ID.String()))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("subject with such name already exists"))

				return
			}

			log.Error("failed to update subject", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update subject"))

			return
		}

		log.Info("subject successfully updated")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, UpdateResponse{
			Response: resp.OK(),
			Subject:  ToModelOne(subject),
		})
	}
}
