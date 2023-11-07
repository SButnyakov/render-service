package auth

import (
	resp "backend-api/internal/lib/api/response"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func New(log *slog.Logger, m *tokenManager.Manager, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				log.Error("empty authorization header")
				responseError(w, r, resp.Error("empty authorization header"), http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				log.Error("invalid authorization header")
				responseError(w, r, resp.Error("invalid authorization header"), http.StatusUnauthorized)
				return
			}

			if len(headerParts[1]) == 0 {
				log.Error("empty authorization token")
				responseError(w, r, resp.Error("empty authorization token"), http.StatusUnauthorized)
				return
			}

			claims, err := m.Parse(headerParts[1], secret)
			if err != nil {
				log.Error("failed to parse token", sl.Err(err))
				responseError(w, r, resp.Error(err.Error()), http.StatusUnauthorized)
				return
			}
			if time.Now().Unix() > claims.ExpiresAt {
				log.Error("expired refresh token")
				responseError(w, r, resp.Error("token has expired"), http.StatusUnauthorized)
				return
			}

			id, err := strconv.Atoi(claims.Subject)
			uid := int64(id)
			if err != nil {
				log.Error("invalid payload", sl.Err(err))
				responseError(w, r, resp.Error("invalid payload"), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "uid", uid)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}
