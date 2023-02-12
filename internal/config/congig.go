package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type Config struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	ServerAddr     string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"` //addr:port
}

func LoadConfig() Config {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

var AppConfig Config = LoadConfig()
