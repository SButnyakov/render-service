package main

import (
	"backend-api/internal/config"
	"backend-api/internal/lib/logger/sl"
	"backend-api/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Envs
	cfgPath := os.Getenv("BUFFER_CONFIG_PATH")
	// inputPath := os.Getenv("BUFFER_INPUT_PATH")
	// outputPath := os.Getenv("BUFFER_OUTPUT_PATH")

	// Config
	cfg := config.MustLoad(cfgPath)

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log = log.With("env", cfg.Env)
	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))
	log.Debug("logger debug mode enabled")

	client, err := redis.New(cfg.Redis.Address)
	if err != nil {
		log.Error("failed to initialize redis", sl.Err(err))
		os.Exit(-1)
	}
	defer client.Close()

	router := chi.NewRouter()

	// Server
	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
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

/*
f := func() {
		priorityInARowCounter := 0
		for {
			var data []string
			if priorityInARowCounter < cfg.Redis.MaxPriorityInARow {
				fmt.Print("checking priority... ")
				data, err = client.BLPop(context.Background(), 2*time.Second, cfg.Redis.PriorityQueueName).Result()
				priorityInARowCounter++
			} else {
				fmt.Print("checking usual... ")
				data, err = client.BLPop(context.Background(), 2*time.Second, cfg.Redis.QueueName).Result()
				priorityInARowCounter = 0
			}

			if err != nil {
				if errors.Is(err, redis.Nil) {
					fmt.Println("nothing")
					continue
				}
				log.Error("reading redis fail", sl.Err(err))
			}

			var newOrder storage.RedisData

			b := []byte(data[1])
			err = json.Unmarshal(b, &newOrder)
			if err != nil {
				log.Error("failed to unmarshal new order", sl.Err(err))
			}
			fmt.Println(newOrder)
		}
	}

	go f()
*/
