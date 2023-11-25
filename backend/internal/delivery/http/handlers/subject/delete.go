package subject

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	"ufahack_2023/internal/delivery/http/handlers/common"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/storage"
	resp "ufahack_2023/pkg/api/response"
	"ufahack_2023/pkg/logger/sl"
)

type Deleter interface {
	DeleteSubject(ctx context.Context, id domain.ID) error
}

type DeleteResponse struct {
	resp.Response
	SubjectID domain.ID
}

func NewDeleteSubject(
	log *slog.Logger,
	deleter Deleter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.subject.DeleteSubjects"

		log = log.With(
			sl.Op(op),
		)

		uid := common.ExtractUUIDParam(log, w, r, "subjectID")
		if uid == uuid.Nil {
			return
		}

		err := deleter.DeleteSubject(r.Context(), uid)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Warn("subject not found")

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("subject not found"))

				return
			}

			log.Error("failed to delete subject", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete subject"))

			return
		}

		log.Info("subject successfully deleted", slog.String("uid", uid.String()))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, DeleteResponse{
			Response:  resp.OK(),
			SubjectID: uid,
		})
	}
}
