package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/server/api/handlers/subscribe"
	"backend-api/internal/server/middleware/cors"
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

func main() {
	// Envs
	cfgPath := os.Getenv("API_CONFIG_PATH")
	storagePath := os.Getenv("API_STORAGE_PATH")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	// DB
	pg, err := postgres.New(storagePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Db.Close()

	// Repos
	orders := repos.NewOrderRepository(pg)
	payments := repos.NewPaymentRepository(pg)
	_ = orders
	_ = payments

	// JWT manager

	// Router
	router := chi.NewRouter()

	// Router middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.New())

	// Router handlers
	router.Post("/subscribe", subscribe.New(log))

	// Server
	server := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Startup
	log.Info("starting server", slog.String("address", cfg.Address))
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(-1)
	}
}
