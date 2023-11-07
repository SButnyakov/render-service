package edit

import (
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

const (
	packagePath = "server.auth.handlers.edit."
)

type Request struct {
	Login    string `json:"login" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type Response struct {
	resp.Response
}

type UserUpdater interface {
	Update(uid int64, user storage.User) error
}

func New(log *slog.Logger, updater UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			responseError(w, r, resp.Error("empty request"), http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			responseError(w, r, resp.Error("failed to decode request"), http.StatusBadRequest)
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			responseError(w, r, resp.ValidationError(validateErr), http.StatusBadRequest)
			return
		}

		user := storage.User{Login: req.Login, Email: req.Email, Password: req.Password}
		err = updater.Update(r.Context().Value("uid").(int64), user)
		if err != nil {
			log.Error("failed to update user info", sl.Err(err))
			responseError(w, r, resp.Error(err.Error()), http.StatusBadRequest)
		}

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
