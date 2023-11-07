package update

import (
	"log/slog"
	"net/http"
)

const (
	packagePath = "server.auth.handlers.update."
)

type Request struct {
	Login    string `json:"login" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type UserUpdater interface {
	Update(uid int64) error
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "New"

		w.WriteHeader(200)
	}
}
