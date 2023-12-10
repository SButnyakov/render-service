package blend

import (
	resp "backend-api/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type OrderStatusUpdater interface {
	UpdateStatus(string, int64, int64) error
}

type Response struct {
	resp.Response
}

func NewUpdate(log *slog.Logger, updater OrderStatusUpdater, orderStatusesMap map[string]int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
