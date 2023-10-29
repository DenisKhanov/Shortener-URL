package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		flagArgs       []string
		expectedConfig *ENVConfig
	}{
		{
			name:     "default values",
			flagArgs: []string{},
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:8080",
				EnvBaseURL:     "http://localhost:8080",
				EnvStoragePath: "/tmp/short-url-db.json",
				EnvLogLevel:    "info",
			},
		},
		{
			name: "environment variables",
			envVars: map[string]string{
				"SERVER_ADDRESS":    "localhost:9090",
				"BASE_URL":          "http://localhost:9090",
				"FILE_STORAGE_PATH": "/tmp/test.json",
				"LOG_LEVEL":         "debug",
			},
			flagArgs: []string{},
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:9090",
				EnvBaseURL:     "http://localhost:9090",
				EnvStoragePath: "/tmp/test.json",
				EnvLogLevel:    "debug",
			},
		},
		{
			name: "command line flags",
			flagArgs: []string{
				"-a", "localhost:7070",
				"-b", "http://localhost:7070",
				"-f", "/tmp/flag-test.json",
				"-l", "error",
			},
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:7070",
				EnvBaseURL:     "http://localhost:7070",
				EnvStoragePath: "/tmp/flag-test.json",
				EnvLogLevel:    "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup and clear flags and environment variables
			oldArgs := os.Args
			oldCommandLine := flag.CommandLine
			defer func() {
				os.Args = oldArgs
				flag.CommandLine = oldCommandLine
			}()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key) // clean up
			}

			// Set command line flags
			os.Args = append([]string{"cmd"}, tt.flagArgs...)

			// Call the function under test
			gotConfig := NewConfig()

			// Assert the result
			assert.Equal(t, tt.expectedConfig, gotConfig)
		})
	}
}

//
//func TestNewConfig(t *testing.T) {
//	tests := []struct {
//		name           string
//		envServAddr    string
//		envBaseURL     string
//		envStoragePath string
//		//envLogLevel string
//		args     []string
//		expected *ENVConfig
//	}{
//		{
//			name:     "test config not environment & not flags",
//			args:     []string{"cmd"},
//			expected: &ENVConfig{EnvServAdr: "localhost:8080", EnvBaseURL: "http://localhost:8080", EnvStoragePath: "/tmp/short-url-db.json"},
//		},
//		{
//			name:     "test config not environment",
//			args:     []string{"cmd", "-a", "localhost:9090", "-b", "http://flags", "-f", "/tmp/flag.json"},
//			expected: &ENVConfig{EnvServAdr: "localhost:9090", EnvBaseURL: "http://flags", EnvStoragePath: "/tmp/flag.json"},
//		},
//		{
//			name:     "test config flag -a not environment",
//			args:     []string{"cmd", "-a", "localhost:9090"},
//			expected: &ENVConfig{EnvServAdr: "localhost:9090", EnvBaseURL: "http://localhost:8080", EnvStoragePath: "/tmp/short-url-db.json"},
//		},
//		{
//			name:           "test config environment & flags",
//			envServAddr:    "localhost:9090",
//			envBaseURL:     "http://enviroment",
//			envStoragePath: "/tmp/env.json",
//			//envLogLevel: "warn",
//			args:     []string{"cmd", "-a", "localhost:7070", "-b", "http://flags", "-f", "/tmp/flag.json"},
//			expected: &ENVConfig{EnvServAdr: "localhost:9090", EnvBaseURL: "http://enviroment", EnvStoragePath: "/tmp/env.json"},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.envServAddr != "" {
//				os.Setenv("SERVER_ADDRESS", tt.envServAddr)
//			}
//			if tt.envBaseURL != "" {
//				os.Setenv("BASE_URL", tt.envBaseURL)
//			}
//			if tt.envStoragePath != "" {
//				os.Setenv("FILE_STORAGE_PATH", tt.envStoragePath)
//			}
//
//			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Сбрасываем значение флагов перед каждым тестом
//			os.Args = tt.args
//			cfg := NewConfig()
//			assert.Equal(t, tt.expected, cfg)
//		})
//	}
//}
