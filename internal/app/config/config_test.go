package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		envServAddr string
		envBaseURL  string
		args        []string
		expected    *ENVConfig
	}{
		// ... existing tests
		{
			args:     []string{"cmd"},
			expected: &ENVConfig{EnvServAdr: "localhost:8080", EnvBaseURL: "http://localhost:8080"},
		},
		{
			envServAddr: "localhost:9090",
			envBaseURL:  "http://example.com",
			args:        []string{"cmd", "-a", "localhost:7070", "-b", "http://test.com"},
			expected:    &ENVConfig{EnvServAdr: "localhost:7070", EnvBaseURL: "http://test.com"},
		},
	}

	for _, test := range tests {
		if test.envServAddr != "" {
			os.Setenv("SERVER_ADDRESS", test.envServAddr)
		} else {
			os.Unsetenv("SERVER_ADDRESS")
		}
		if test.envBaseURL != "" {
			os.Setenv("BASE_URL", test.envBaseURL)
		} else {
			os.Unsetenv("BASE_URL")
		}

		oldArgs := os.Args
		os.Args = test.args
		defer func() { os.Args = oldArgs }()

		cfg := NewConfig()

		assert.Equal(t, test.expected, cfg)
	}
}
