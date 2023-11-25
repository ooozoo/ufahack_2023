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

type Provider interface {
	ListSubjects(ctx context.Context) ([]*domain.Subject, error)
	GetSubjectByID(ctx context.Context, id domain.ID) (*domain.Subject, error)
}

type GetListResponse struct {
	resp.Response
	Subjects []*Subject `json:"subject"`
}

func NewListSubjects(
	log *slog.Logger,
	provider Provider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.subject.ListSubjects"

		log = log.With(
			sl.Op(op),
		)

		subjects, err := provider.ListSubjects(r.Context())
		if err != nil {
			log.Error("failed to get subject", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get subject"))

			return
		}

		log.Debug("successfully fetched subject", slog.Any("subject", subjects))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetListResponse{
			Response: resp.OK(),
			Subjects: ToModelMany(subjects),
		})
	}
}

type GetResponse struct {
	resp.Response
	Subject *Subject `json:"subject"`
}

func NewGetSubject(
	log *slog.Logger,
	provider Provider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.subject.GetSubject"

		log = log.With(
			sl.Op(op),
		)

		uid := common.ExtractUUIDParam(log, w, r, "subjectID")
		if uid == uuid.Nil {
			return
		}

		subject, err := provider.GetSubjectByID(r.Context(), uid)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Warn("subject not found", slog.Any("uid", uid))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("subject not found"))
				return
			}

			log.Error("failed to get subject", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get subject"))
			return
		}

		log.Debug("subject successfully fetched", slog.Any("subject", subject))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetResponse{
			Response: resp.OK(),
			Subject:  ToModelOne(subject),
		})
	}
}
