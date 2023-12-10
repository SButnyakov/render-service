package request

import (
	"backend-api/internal/config"
	resp "backend-api/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

const (
	packagePath = "server.buffer.handler.request."
)

type Response struct {
	resp.Response
	Format       string `json:"format"`
	Resolution   string `json:"resolution"`
	DownloadLink string `json:"download_link"`
}

func New(log *slog.Logger, client *redis.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
