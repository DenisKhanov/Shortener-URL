// Package config provides functionality for configuring the application using environment variables.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"os"
)

// ENVConfig holds configuration settings extracted from environment variables.
// This struct is used to configure various aspects of the application.
type ENVConfig struct {
	ConfigFile     string `env:"CONFIG"`
	EnvServAdr     string `env:"SERVER_ADDRESS"`
	EnvBaseURL     string `env:"BASE_URL"`
	EnvStoragePath string `env:"FILE_STORAGE_PATH"`
	EnvLogLevel    string `env:"LOG_LEVEL"`
	EnvDataBase    string `env:"DATABASE_DSN"`
	EnvHTTPS       string `env:"ENABLE_HTTPS"`
}

func checkConfigFile() *ENVConfig {
	var cfg ENVConfig

	// Parse command line flags
	flag.StringVar(&cfg.ConfigFile, "c", "", "Path to the configuration file")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	if cfg.ConfigFile != "" {
		if err = readConfigFile(&cfg, cfg.ConfigFile); err != nil {
			logrus.Fatal(err)
		}
	}
	// Parse config from JSON file if provided
	if cfgFile := getConfigFilePath(); cfgFile != "" {
		if err = readConfigFile(&cfg, cfgFile); err != nil {
			logrus.Fatal(err)
		}
	}

	return &cfg
}

// NewConfig creates a new ENVConfig instance by parsing command line flags and environment variables.
func NewConfig() *ENVConfig {
	var cfg ENVConfig

	checkConfigFile()

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

// getConfigFilePath returns the path to the config file specified by the -c flag or the CONFIG environment variable.
func getConfigFilePath() string {
	cfgFile := os.Getenv("CONFIG")
	if cfgFile != "" {
		return cfgFile
	}
	return ""
}

func readConfigFile(cfg *ENVConfig, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	return nil
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
