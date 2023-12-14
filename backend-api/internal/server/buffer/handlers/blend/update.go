package blend

import (
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	PackagePath = "server.buffer.handlers.blend.update."
)

type OrderStatusUpdater interface {
	UpdateStatus(string, int64, int64) error
}

type Response struct {
	resp.Response
}

func NewUpdate(log *slog.Logger, updater OrderStatusUpdater, orderStatusesMap map[string]int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = PackagePath + "NewUpdate"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		uid := chi.URLParam(r, "uid")
		fileName := chi.URLParam(r, "filename")
		status := strings.ReplaceAll(chi.URLParam(r, "status"), "-", " ")

		statusId, ok := orderStatusesMap["status"]
		if !ok {
			log.Error("status " + status + " not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(uid)
		if err != nil {
			log.Error("invalid user id", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updater.UpdateStatus(fileName, int64(id), statusId)
		if err != nil {
			log.Error("failed to update status", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
