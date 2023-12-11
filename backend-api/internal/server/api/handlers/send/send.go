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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	PackagePath = "server.buffer.handlers.send."
)

type OrderProvider interface {
	Create(storage.Order) error
}

type Response struct {
	resp.Response
}

func New(log *slog.Logger, inputPath string, cfg *config.Config, provider OrderProvider,
	orderStatuses map[string]int64, redis *redis.Client) http.HandlerFunc {
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

		format := r.FormValue("format")
		resolution := r.FormValue("resolution")

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error("failed to get file from form", sl.Err(err))
			responseError(w, r, resp.Error("failed to get file from form"), http.StatusBadRequest)
			return
		}
		defer file.Close()

		uid := r.Context().Value("uid").(int64)

		userPath := filepath.Join(inputPath, strconv.FormatInt(uid, 10))
		if err = os.MkdirAll(userPath, os.ModePerm); err != nil {
			log.Error("failed to create user's dir", sl.Err(err))
			responseError(w, r, resp.Error("failed to save file"), http.StatusInternalServerError)
			return
		}

		storingName := strconv.FormatInt(time.Now().Unix(), 10) + "." + strings.Split(handler.Filename, ".")[1]

		savePath := userPath + "/" + storingName

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
			os.Remove(savePath)
			return
		}

		err = provider.Create(storage.Order{
			FileName:     handler.Filename,
			StoringName:  storingName,
			CreationDate: time.Now(),
			UserId:       uid,
			StatusId:     orderStatuses["in queue"],
		})
		if err != nil {
			log.Error("failed to create order record", sl.Err(err))
			responseError(w, r, resp.Error("failed to create new order"), http.StatusInternalServerError)
			os.Remove(savePath)
			return
		}

		b, err := json.Marshal(storage.RedisData{
			Format:     format,
			Resolution: resolution,
			SavePath:   savePath,
		})
		if err != nil {
			log.Error("failed to marshal new order", sl.Err(err))
			responseError(w, r, resp.Error("failed to save file"), http.StatusInternalServerError)
			return
		}

		redis.RPush(context.Background(), cfg.Redis.QueueName, string(b))

		responseOK(w, r)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
