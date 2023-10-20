package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type ENVConfig struct {
	EnvServAdr string `env:"SERVER_ADDRESS"`
	EnvBaseURL string `env:"BASE_URL"`
	//EnvLogLevel string `env:"LOG_LEVEL"`
}

func NewConfig() *ENVConfig {
	var cfg ENVConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.EnvServAdr == "" || cfg.EnvBaseURL == "" {
		if cfg.EnvServAdr == "" {
			flag.StringVar(&cfg.EnvServAdr, "a", "localhost:8080", "HTTP server address")
		}
		if cfg.EnvBaseURL == "" {
			flag.StringVar(&cfg.EnvBaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
		}
		//if cfg.EnvLogLevel == "" {
		//	flag.StringVar(&cfg.EnvLogLevel, "l", "info", "Logs level")
		//}
		flag.Parse()
	}

	return &cfg
}
