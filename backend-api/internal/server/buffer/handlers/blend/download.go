package blend

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

const (
	packagePath = "server.buffer.handlers.download."
)

func NewDownload(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
