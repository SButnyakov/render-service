package image

import (
	"backend-api/internal/lib/logger/sl"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	PackagePath = "server.buffer.handlers.image.upload."
)

type OrdersProvider interface {
	UpdateStatus(storingName string, uid, statusId int64) error
	UpdateDownloadLink(uid int64, storingName, downloadLink string) error
}

func NewUpload(log *slog.Logger, inputPath, outputPath string, provider OrdersProvider,
	orderStatusesMap map[string]int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = PackagePath + "NewUpload"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed to parse form", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error("failed to get file from form", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		uid := chi.URLParam(r, "uid")

		userPath := filepath.Join(outputPath, uid)
		if err = os.MkdirAll(userPath, os.ModePerm); err != nil {
			log.Error("failed to create user's dir", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		savePath := filepath.Join(userPath, handler.Filename)

		f, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Error("failed to save file", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			log.Error("failed to save file", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			os.Remove(savePath)
			return
		}

		id, err := strconv.Atoi(uid)
		if err != nil {
			log.Error("failed to cast uid to int64", sl.Err(err))
		}

		err = provider.UpdateStatus(strings.Split(handler.Filename, ".")[0], int64(id), orderStatusesMap["success"])
		if err != nil {
			log.Error("failed to update order status", sl.Err(err))
		}

		downloadLink := fmt.Sprintf("http://%s/%s/images/download/%s", "localhost:8082", uid, handler.Filename)

		err = provider.UpdateDownloadLink(int64(id), strings.Split(handler.Filename, ".")[0], downloadLink)
		if err != nil {
			log.Error("failed to update download link", sl.Err(err))
		}

		err = os.Remove(filepath.Join(inputPath, uid, strings.Split(handler.Filename, ".")[0]))
		if err != nil {
			log.Error("failed to delete source file", sl.Err(err))
		}

		w.WriteHeader(http.StatusOK)
	}
}
