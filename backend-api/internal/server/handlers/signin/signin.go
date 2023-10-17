package signin

import (
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	LoginOrEmail string `json:"login_or_email" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type UserProvider interface {
	CheckCredentials(string, string) (int64, error)
}

func New(log *slog.Logger, userProvider UserProvider, m *tokenManager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.signin.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			responseError(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			responseError(w, r, resp.Error("failed to decode request"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			responseError(w, r, resp.ValidationError(validateErr))
		}

		id, err := userProvider.CheckCredentials(req.LoginOrEmail, req.Password)
		if errors.Is(err, storage.ErrInvalidCredentials) {
			responseError(w, r, resp.Error("invalid credentials"))
			return
		}
		if err != nil {
			responseError(w, r, resp.Error("server-side authorization failed"))
			return
		}

		log.Info("creating new jwt token", slog.Int64("id", id))

		token, err := m.NewJWT(id, time.Hour*72)
		if err != nil {
			responseError(w, r, resp.Error("failed to create jwt token"))
			return
		}

		responseOK(w, r, token)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response) {
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
