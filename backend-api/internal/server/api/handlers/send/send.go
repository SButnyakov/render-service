package send

import (
	"backend-api/internal/config"
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/storage"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const (
	PackagePath = "server.buffer.handlers.send."
)

func New(log *slog.Logger, inputPath string, redis *redis.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = PackagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed to parse form", sl.Err(err))
			responseError(w, r, resp.Error("failed to parse form"), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error("failed to get file from form", sl.Err(err))
			responseError(w, r, resp.Error("failed to get file from form"), http.StatusBadRequest)
			return
		}
		defer file.Close()

		savePath := inputPath + handler.Filename

		f, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Error("failed to save file", sl.Err(err))
			responseError(w, r, resp.Error("failed to save file"), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			log.Error("failed to save file", sl.Err(err))
			responseError(w, r, resp.Error("failed to save file"), http.StatusInternalServerError)
			return
		}

		newOrder := storage.RedisData{
			SavePath: savePath,
		}
		b, err := json.Marshal(newOrder)
		if err != nil {
			log.Error("failed to save file", sl.Err(err))
			responseError(w, r, resp.Error("failed to save file"), http.StatusInternalServerError)
			return
		}

		redis.RPush(context.Background(), cfg.Redis.QueueName, string(b))

		w.WriteHeader(http.StatusOK)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
