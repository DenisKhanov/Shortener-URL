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
	EnvTLS         string `env:"ENABLE_TLS"`
	EnvSubnet      string `env:"TRUSTED_SUBNET"`
	EnvGRPC        string `env:"GRPC_SERVER"`
}

// NewConfig creates a new ENVConfig instance by parsing command line flags and environment variables.
func NewConfig() (*ENVConfig, error) {
	var cfg ENVConfig

	flag.StringVar(&cfg.ConfigFile, "c", "", "Enter path to config file Or use CONFIG env")

	flag.StringVar(&cfg.EnvServAdr, "a", "localhost:8080", "Enter URLServer as ip_address:port or use SERVER_ADDRESS env")

	flag.StringVar(&cfg.EnvBaseURL, "b", "http://localhost:8080", "Enter URLPrefix as http://ip_address:port or use BASE_URL env")

	flag.StringVar(&cfg.EnvStoragePath, "f", "/tmp/short-url-db.json", "Enter path for saving data file or use FILE_STORAGE_PATH env")

	flag.StringVar(&cfg.EnvLogLevel, "l", "info", "Enter logg level or use LOG_LEVEL env")

	flag.StringVar(&cfg.EnvDataBase, "d", "", "Enter url to connect database as host=host port=port user=postgres password=postgres "+
		"dbname=dbname sslmode=disable or use DATABASE_DSN env")

	flag.StringVar(&cfg.EnvTLS, "s", "", "Enter TLS on enable or disable")

	flag.StringVar(&cfg.EnvSubnet, "t", "", "Enter trusted subnet or use TRUSTED_SUBNET env")

	flag.StringVar(&cfg.EnvGRPC, "g", ":3200", "Enter gRPC server address or use GRPC_SERVER env")

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
		cfg1.EnvTLS = cfgFromFile.EnvTLS
	}
	if flag.Lookup("t") == nil {
		cfg1.EnvSubnet = cfgFromFile.EnvSubnet
	}
	if flag.Lookup("g") == nil {
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
	fmt.Printf("Build version: %s\n", getValueOrDefault(buildVersion, "2.0"))
	fmt.Printf("Build date: %s\n", getValueOrDefault(buildDate, "13.04.2024"))
	fmt.Printf("Build commit: %s\n", getValueOrDefault(buildCommit, "Clean architect"))
}
