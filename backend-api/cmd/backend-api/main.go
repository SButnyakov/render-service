package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/server/handlers/signin"
	"backend-api/internal/server/handlers/signup"
	mwLogger "backend-api/internal/server/middleware/logger"
	"backend-api/internal/storage/postgres"
	"backend-api/internal/storage/postgres/repos"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	pg, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Db.Close()

	users := repos.NewUserStorage(pg)
	order := repos.NewOrderStorage(pg)
	_ = order

	m, err := tokenManager.New("secret-key")
	if err != nil {
		log.Error("failed to initialie jwt manager", sl.Err(err))
		os.Exit(-1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/signup", signup.New(log, users))
	router.Post("/signin", signin.New(log, users, m))

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", cfg.Address))
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
