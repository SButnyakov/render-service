package signin

import (
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/api/tokenManager"
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
	LoginOrEmail string `json:"login_or_email" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProvider interface {
	CheckCredentials(string, string) (int64, error)
	UpdateRefreshToken(int64, string) error
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

		uid, err := userProvider.CheckCredentials(req.LoginOrEmail, req.Password)
		if errors.Is(err, storage.ErrInvalidCredentials) {
			log.Error("invalid credentials", sl.Err(err))
			responseError(w, r, resp.Error("invalid credentials"), http.StatusBadRequest)
			return
		}
		if err != nil {
			responseError(w, r, resp.Error("server-side authorization failed"), http.StatusInternalServerError)
			return
		}

		log.Info("creating new jwt token", slog.Int64("uid", uid))

		// TODO: implement generating new tokens

		accessToken, err := m.NewJWT(uid)
		if err != nil {
			log.Error("failed to generate a new access token", sl.Err(err))
			responseError(w, r, resp.Error("failed to create jwt token"), http.StatusInternalServerError)
			return
		}

		refreshToken, err := m.NewRT(uid)
		if err != nil {
			log.Error("failed to generate a new refresh token", sl.Err(err))
			responseError(w, r, resp.Error("failed to create jwt token"), http.StatusInternalServerError)
			return
		}

		err = userProvider.UpdateRefreshToken(uid, refreshToken)
		if err != nil {
			log.Error("failed to set refresh token", sl.Err(err))
			responseError(w, r, resp.Error("failed to set refresh token"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, accessToken, refreshToken)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request, accessToken, refreshToken string) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
