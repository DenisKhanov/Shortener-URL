package config

import "flag"

type Config struct {
	FlagRunAddr string
	BaseURL     string
}

func NewConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.FlagRunAddr, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.Parse()
	return cfg
}
