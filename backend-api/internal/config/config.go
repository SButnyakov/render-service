package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env           string `yaml:"env" env-default:"prod"`
	HTTPServer    `yaml:"http_server"`
	DB            `yaml:"db"`
	Redis         `yaml:"redis"`
	Payments      `yaml:"payments"`
	Subscriptions `yaml:"subscriptions"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Name     string `yaml:"name" env-required:"true"`
	User     string `yaml:"user" env-default:"user"`
	Password string `yaml:"password" env-default:"password"`
	Host     string `yaml:"host" env-default:"postgres"`
	Port     string `yaml:"port" env-default:"5432"`
}

type Redis struct {
	Address           string `yaml:"address" env-default:"localhost:6379"`
	QueueName         string `yaml:"queue_name" env-default:"render-list"`
	PriorityQueueName string `yaml:"priority_queue_name" env-default:"render-list"`
	MaxPriorityInARow int    `yaml:"max_priority_in_a_row" env-default:"5"`
}

type Payments struct {
	SubPremiumMonth string `yaml:"sub_premium_month" env-default:"sub-premium-month"`
}

type Subscriptions struct {
	Premium string `yaml:"premium" env-default:"premium"`
}

func MustLoad(configPath string) *Config {
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
