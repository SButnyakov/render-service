package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"prod"`
	HTTPServer `yaml:"http_server"`
	Envs       `yaml:"envs"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-defult:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Envs struct {
	Storage   string `yaml:"storage" env-required:"true"`
	JWTSecret string `yaml:"jwt_secret" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
