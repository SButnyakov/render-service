package refresh

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
	"strconv"
)

const (
	PackagePath = "server.handlers.signin."
)

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenProvider interface {
	GetRefreshToken(uid int64) (string, error)
	UpdateRefreshToken(uid int64, refreshToken string) error
}

func New(log *slog.Logger, provider TokenProvider, m *tokenManager.Manager, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = PackagePath + "New"

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

		claims, err := m.Parse(req.RefreshToken, secret)
		if err != nil {
			log.Error("invalid refresh token", sl.Err(err))
			responseError(w, r, resp.Error("invalid refresh token"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(claims.Subject)
		uid := int64(id)
		if err != nil {
			log.Error("invalid payload", sl.Err(err))
			responseError(w, r, resp.Error("invalid payload"), http.StatusBadRequest)
			return
		}

		token, err := provider.GetRefreshToken(uid)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("no user with this refresh token")
				responseError(w, r, resp.Error("refresh token was not found"), http.StatusUnauthorized)
				return
			}

			log.Error("failed to get refresh token", sl.Err(err))
			responseError(w, r, resp.Error("failed to get refresh token"), http.StatusInternalServerError)
			return
		}

		if token != req.RefreshToken {
			log.Error("refresh tokens do not match")
			responseError(w, r, resp.Error("failed to get refresh token"), http.StatusUnauthorized)
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

		err = provider.UpdateRefreshToken(uid, req.RefreshToken)
		if err != nil {
			log.Error("failed to update refresh token")
			responseError(w, r, resp.Error("failed to update refresh token"), http.StatusInternalServerError)
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
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
