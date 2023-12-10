package image

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func NewUpload(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "NewUpload"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
	}
}
