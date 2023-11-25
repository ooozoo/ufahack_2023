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

type Saver interface {
	SaveSubject(ctx context.Context, name string) (domain.ID, error)
}

type CreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateResponse struct {
	resp.Response
	SubjectID domain.ID `json:"subject_id"`
}

func NewCreateSubject(
	log *slog.Logger,
	saver Saver,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.subject.CreateSubject"

		log = log.With(
			sl.Op(op),
		)

		var req CreateRequest

		if !common.DecodeRequest(log, w, r, &req) {
			return
		}
		if !common.ValidateRequest(log, w, r, req) {
			return
		}

		uid, err := saver.SaveSubject(r.Context(), req.Name)
		if err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				log.Warn("subject already exists", slog.String("name", req.Name))

				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("subject already exists"))

				return
			}

			log.Error("failed to create subject", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to create subject"))

			return
		}

		log.Info("created subject", slog.String("id", uid.String()))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, CreateResponse{
			Response:  resp.OK(),
			SubjectID: uid,
		})
	}
}
