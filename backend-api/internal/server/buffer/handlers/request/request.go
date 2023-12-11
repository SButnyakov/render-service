package request

import (
	"backend-api/internal/config"
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	packagePath = "server.buffer.handler.request."
)

type Response struct {
	resp.Response
	Format       string `json:"format,omitempty"`
	Resolution   string `json:"resolution,omitempty"`
	DownloadLink string `json:"download_link,omitempty"`
}

func New(log *slog.Logger, client *redis.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		data, err := client.BLPop(context.Background(), 1*time.Second, cfg.Redis.PriorityQueueName).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Error("reading redis priority queue fail", sl.Err(err))
		}
		if err != nil {
			data, err = client.BLPop(context.Background(), 1*time.Second, cfg.Redis.QueueName).Result()
		}
		if err != nil {
			if errors.Is(err, redis.Nil) {
				log.Info("empty queue's")
				responseEmpty(w, r)
				return
			}
			log.Error("reading redis queue fail")
			responseError(w, r, resp.Error("reading orders list failed"), http.StatusInternalServerError)
			return
		}

		var newOrder storage.RedisData

		b := []byte(data[1])
		err = json.Unmarshal(b, &newOrder)
		if err != nil {
			log.Error("failed to unmarshal new order", sl.Err(err))
		}

		pathList := strings.Split(newOrder.SavePath, "/")
		listLength := len(pathList)

		downloadLink := fmt.Sprintf("http://%s/%s/blend/download/%s", "localhost:8082", pathList[listLength-2], pathList[listLength-1])

		responseOK(w, r, newOrder.Format, newOrder.Resolution, downloadLink)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request, format, resolution, downloadLink string) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		Format:       format,
		Resolution:   resolution,
		DownloadLink: downloadLink,
	})
}

func responseEmpty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.Empty(),
	})
}
