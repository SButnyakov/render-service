package signup

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

type Request struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type UserCreator interface {
	Create(storage.User) error
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.signup.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		user := storage.User{Login: req.Login, Email: req.Email, Password: req.Password}
		err = userCreator.Create(user)
		if errors.Is(err, storage.ErrUserExists) {
			render.JSON(w, r, resp.Error("user with the same login or email already exists"))
			return
		}
		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
