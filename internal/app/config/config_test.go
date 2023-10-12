package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-a", "localhost:9090", "-b", "http://localhost:8000"}

	cfg := NewConfig()

	// убедитесь, что значения были правильно установлены
	assert.Equal(t, "localhost:9090", cfg.FlagRunAddr)
	assert.Equal(t, "http://localhost:8000", cfg.BaseURL)
}
