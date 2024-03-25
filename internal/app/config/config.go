// Package config provides functionality for configuring the application using environment variables.
package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

// ENVConfig holds configuration settings extracted from environment variables.
// This struct is used to configure various aspects of the application.
type ENVConfig struct {
	EnvServAdr     string `env:"SERVER_ADDRESS"`
	EnvBaseURL     string `env:"BASE_URL"`
	EnvStoragePath string `env:"FILE_STORAGE_PATH"`
	EnvLogLevel    string `env:"LOG_LEVEL"`
	EnvDataBase    string `env:"DATABASE_DSN"`
	EnvHTTPS       string `env:"ENABLE_HTTPS"`
}

// NewConfig creates a new ENVConfig instance by parsing command line flags and environment variables.
func NewConfig() *ENVConfig {
	var cfg ENVConfig

	// Parse command line flags
	flag.StringVar(&cfg.EnvServAdr, "a", "localhost:8080", "HTTP server address")

	flag.StringVar(&cfg.EnvBaseURL, "b", "http://localhost:8080", "Base URL for shortened links")

	flag.StringVar(&cfg.EnvStoragePath, "f", "/tmp/short-url-db.json", "Path for saving data file")

	flag.StringVar(&cfg.EnvLogLevel, "l", "info", "Set logg level")

	flag.StringVar(&cfg.EnvDataBase, "d", "", "Set connect DB config")

	flag.StringVar(&cfg.EnvHTTPS, "s", "", "Set HTTPS on enable")

	flag.Parse()

	// Parse environment variables.
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	return &cfg
}

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// getValueOrDefault returns the value, and if it is empty,it returns the default value.
func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// PrintProjectInfo print info (version,date,commit) about build.
func PrintProjectInfo() {
	fmt.Printf("Build version: %s\n", getValueOrDefault(buildVersion, "N/A"))
	fmt.Printf("Build date: %s\n", getValueOrDefault(buildDate, "N/A"))
	fmt.Printf("Build commit: %s\n", getValueOrDefault(buildCommit, "N/A"))
}
