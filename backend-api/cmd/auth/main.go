package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/server/auth/handlers/edit"
	"backend-api/internal/server/auth/handlers/refresh"
	"backend-api/internal/server/auth/handlers/signin"
	"backend-api/internal/server/auth/handlers/signup"
	"backend-api/internal/server/middleware/auth"
	"backend-api/internal/server/middleware/cors"
	mwLogger "backend-api/internal/server/middleware/logger"
	"backend-api/internal/storage/postgres"
	"backend-api/internal/storage/postgres/repos"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Envs
	cfgPath := os.Getenv("AUTH_CONFIG_PATH")
	storagePath := os.Getenv("AUTH_STORAGE_PATH")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing service", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	// DB
	pg, err := postgres.New(storagePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Db.Close()

	// Repos
	users := repos.NewUserRepository(pg)

	// JWT manager
	jwtManager, err := tokenManager.New(jwtSecretKey)
	if err != nil {
		log.Error("failed to initialie jwt manager", sl.Err(err))
		os.Exit(-1)
	}

	// Router
	router := chi.NewRouter()

	// Router middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.New())

	// Router handlers
	router.Post("/signup", signup.New(log, users))
	router.Post("/signin", signin.New(log, users, jwtManager))
	router.Put("/refresh", refresh.New(log, users, jwtManager, jwtSecretKey))

	router.Group(func(r chi.Router) {
		r.Use(auth.New(log, jwtManager, jwtSecretKey))
		r.Put("/user/edit", edit.New(log, users))
	})

	// Server
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
		os.Exit(-1)
	}
}
