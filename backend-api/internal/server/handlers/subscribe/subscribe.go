package subscribe

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "handlers.subscribe.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
