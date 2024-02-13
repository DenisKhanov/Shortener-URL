package config

import (
	"bytes"
	"flag"
	"fmt"
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
			name: "default values",
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:8080",
				EnvBaseURL:     "http://localhost:8080",
				EnvStoragePath: "/tmp/short-url-db.json",
				EnvLogLevel:    "info",
				EnvDataBase:    "",
			},
		},
		{
			name: "environment variables",
			envVars: map[string]string{
				"SERVER_ADDRESS":    "localhost:9090",
				"BASE_URL":          "http://localhost:9090",
				"FILE_STORAGE_PATH": "/tmp/test.json",
				"LOG_LEVEL":         "debug",
				"DATABASE_DSN":      "environment",
			},
			flagArgs: []string{},
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:9090",
				EnvBaseURL:     "http://localhost:9090",
				EnvStoragePath: "/tmp/test.json",
				EnvLogLevel:    "debug",
				EnvDataBase:    "environment",
			},
		},
		{
			name: "command line flags",
			flagArgs: []string{
				"-a", "localhost:7070",
				"-b", "http://localhost:7070",
				"-f", "/tmp/flag-test.json",
				"-l", "error",
				"-d", "flags",
			},
			expectedConfig: &ENVConfig{
				EnvServAdr:     "localhost:7070",
				EnvBaseURL:     "http://localhost:7070",
				EnvStoragePath: "/tmp/flag-test.json",
				EnvLogLevel:    "error",
				EnvDataBase:    "flags",
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

func TestGetValueOrDefault(t *testing.T) {
	// Test case 1: value is not empty
	result := getValueOrDefault("testValue", "defaultValue")
	expected := "testValue"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Test case 2: value is empty
	result = getValueOrDefault("", "defaultValue")
	expected = "defaultValue"
	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestPrintProjectInfo(t *testing.T) {
	// Redirect stdout to capture printed output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the PrintProjectInfo function
	PrintProjectInfo()

	// Read captured output
	var buf bytes.Buffer
	w.Close()
	buf.ReadFrom(r)
	os.Stdout = old

	// Expected output
	expectedOutput := fmt.Sprintf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		getValueOrDefault(buildVersion, "N/A"),
		getValueOrDefault(buildDate, "N/A"),
		getValueOrDefault(buildCommit, "N/A"))

	// Check if the printed output matches the expected output
	if buf.String() != expectedOutput {
		t.Errorf("Expected output: %s, Got: %s", expectedOutput, buf.String())
	}
}
