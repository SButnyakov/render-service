package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/api/tokenManager"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/server/api/handlers/send"
	"backend-api/internal/server/api/handlers/subscribe"
	"backend-api/internal/server/middleware/auth"
	"backend-api/internal/server/middleware/cors"
	mwLogger "backend-api/internal/server/middleware/logger"
	"backend-api/internal/storage/postgres"
	"backend-api/internal/storage/postgres/repos"
	"backend-api/internal/storage/redis"
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
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	inputPath := os.Getenv("FILES_INPUT_PATH")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
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
	paymentTypes := repos.NewPaymentTypeRepository(pg)
	subscriptionTypes := repos.NewSubscriptionTypeRepository(pg)
	subscriptions := repos.NewSubscriptionRepository(pg)

	// Maps
	orderStatusesMap, err := orderStatuses.GetStatusesMap()
	if err != nil {
		log.Error("failed to get order statuses", sl.Err(err))
		os.Exit(-1)
	}
	paymentTypesMap, err := paymentTypes.GetTypesMap()
	if err != nil {
		log.Error("failed to get payment types", sl.Err(err))
		os.Exit(-1)
	}
	subscriptionTypesMap, err := subscriptionTypes.GetTypesMap()
	if err != nil {
		log.Error("failed to get subscription types", sl.Err(err))
		os.Exit(-1)
	}
	_ = orderStatusesMap

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
	router := chi.NewRouter()

	// Router middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(cors.New())
	router.Use(auth.New(log, jwtManager))

	// Router handlers
	router.Post("/subscribe", subscribe.New(log, cfg, paymentTypesMap, subscriptionTypesMap, subscriptions))
	router.Post("/send", send.New(log, inputPath, cfg, orders, orderStatusesMap, client))

	// Server
	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Startup
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(-1)
	}
}
