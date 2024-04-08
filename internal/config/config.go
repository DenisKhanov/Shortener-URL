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
	EnvSubnet      string `env:"TRUSTED_SUBNET"`
}

// NewConfig creates a new ENVConfig instance by parsing command line flags and environment variables.
func NewConfig() (*ENVConfig, error) {
	var cfg ENVConfig

	flag.StringVar(&cfg.ConfigFile, "c", "", "Path to the configuration file")

	flag.StringVar(&cfg.EnvServAdr, "a", "localhost:8080", "HTTP server address")

	flag.StringVar(&cfg.EnvBaseURL, "b", "http://localhost:8080", "Base URL for shortened links")

	flag.StringVar(&cfg.EnvStoragePath, "f", "/tmp/short-url-db.json", "Path for saving data file")

	flag.StringVar(&cfg.EnvLogLevel, "l", "info", "Set logg level")

	flag.StringVar(&cfg.EnvDataBase, "d", "", "Set connect DB config")

	flag.StringVar(&cfg.EnvHTTPS, "s", "", "Set HTTPS on enable")

	flag.StringVar(&cfg.EnvSubnet, "t", "", "Use trusted subnet")

	flag.Parse()

	// Parse config from JSON file if provided
	if cfgFilePath := getConfigFilePath(); cfgFilePath != "" {
		if err := setConfigFromFile(cfgFilePath, &cfg); err != nil {
			logrus.Error(err)
			return nil, err
		}
	} else if cfg.ConfigFile != "" {
		if err := setConfigFromFile(cfg.ConfigFile, &cfg); err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	// Parse environment variables.
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &cfg, nil
}

// getConfigFilePath returns the path to the config file specified by the -c flag or the CONFIG environment variable.
func getConfigFilePath() string {
	cfgFile := os.Getenv("CONFIG")
	return cfgFile
}

func setConfigFromFile(path string, cfg1 *ENVConfig) error {
	var cfgFromFile ENVConfig

	data, err := os.ReadFile(path)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = json.Unmarshal(data, &cfgFromFile)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if flag.Lookup("a") == nil {
		cfg1.EnvServAdr = cfgFromFile.EnvServAdr
	}

	if flag.Lookup("b") == nil {
		cfg1.EnvBaseURL = cfgFromFile.EnvBaseURL
	}

	if flag.Lookup("f") == nil {
		cfg1.EnvStoragePath = cfgFromFile.EnvStoragePath
	}

	if flag.Lookup("l") == nil {
		cfg1.EnvLogLevel = cfgFromFile.EnvLogLevel
	}

	if flag.Lookup("d") == nil {
		cfg1.EnvDataBase = cfgFromFile.EnvDataBase
	}

	if flag.Lookup("s") == nil {
		cfg1.EnvHTTPS = cfgFromFile.EnvHTTPS
	}
	if flag.Lookup("t") == nil {
		cfg1.EnvSubnet = cfgFromFile.EnvSubnet
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
