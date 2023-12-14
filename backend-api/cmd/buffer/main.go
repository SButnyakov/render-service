package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/server/buffer/handlers/blend"
	"backend-api/internal/server/buffer/handlers/download"
	"backend-api/internal/server/buffer/handlers/image"
	"backend-api/internal/server/buffer/handlers/request"
	"backend-api/internal/server/middleware/auth"
	mwLogger "backend-api/internal/server/middleware/logger"
	"backend-api/internal/storage/postgres"
	"backend-api/internal/storage/postgres/repos"
	"backend-api/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Envs
	cfgPath := os.Getenv("BUFFER_CONFIG_PATH")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	inputPath := os.Getenv("FILES_INPUT_PATH")
	outputPath := os.Getenv("FILES_OUTPUT_PATH")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With("env", cfg.Env)
	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))
	log.Debug("logger debug mode enabled")

	// DB
	pg, err := postgres.New(cfg)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Db.Close()

	// Repos
	orderStatuses := repos.NewOrderStatusesRepository(pg)
	orders := repos.NewOrderRepository(pg)

	// Maps
	orderStatusesMap, err := orderStatuses.GetStatusesMap()
	if err != nil {
		log.Error("failed to get order statuses", sl.Err(err))
		os.Exit(-1)
	}

	// Redis
	client, err := redis.New(cfg)
	if err != nil {
		log.Error("failed to initialize redis", sl.Err(err))
		os.Exit(-1)
	}
	defer client.Close()

	// JWT manager
	jwtManager, err := tokenManager.New(jwtSecretKey)
	if err != nil {
		log.Error("failed to initialize jwt token manager", sl.Err(err))
		os.Exit(-1)
	}

	// Router
	mainRouter := chi.NewRouter()

	// Router middleware
	mainRouter.Use(middleware.RequestID)
	mainRouter.Use(mwLogger.New(log))
	mainRouter.Use(middleware.Recoverer)
	mainRouter.Use(middleware.URLFormat)

	// Router handlers
	mainRouter.Get("/request", request.New(log, client, cfg))
	mainRouter.Route("/{uid}", func(uidRouter chi.Router) {
		uidRouter.Route("/blend", func(blendRouter chi.Router) {
			blendRouter.Get("/download/{filename}", download.New(log, inputPath))
			blendRouter.Put("/update/{filename}/{status}", blend.NewUpdate(log, orders, orderStatusesMap))
		})
		uidRouter.Route("/image", func(imageRouter chi.Router) {
			imageRouter.Post("/upload", image.NewUpload(log, inputPath, outputPath, orders, orderStatusesMap))
			imageRouter.Route("/", func(authImageRouter chi.Router) {
				authImageRouter.Use(auth.New(log, jwtManager))
				authImageRouter.Get("/download/{filename}", download.New(log, outputPath))
			})
		})
	})

	// Server
	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mainRouter,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(-1)
	}
}
