package logcfg

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestRunLoggerConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		envLogs     string
		expectedLog *logrus.Logger
	}{
		{
			name:    "Test with 'debug' log level",
			envLogs: "debug",
			expectedLog: &logrus.Logger{
				Out: os.Stdout,
				Formatter: &logrus.TextFormatter{
					CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
						_, filename := path.Split(f.File)
						filename = fmt.Sprintf("%s.%d.%s", filename, f.Line, f.Function)
						return "", filename
					},
				},
				Hooks: make(logrus.LevelHooks),
				Level: logrus.DebugLevel,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the log output
			var logOutput bytes.Buffer
			logrus.SetOutput(&logOutput)

			// Run the logger configuration
			RunLoggerConfig(tt.envLogs)

			// Check log level
			assert.Equal(t, tt.expectedLog.Level, logrus.GetLevel())

			// Check log formatter
			assert.NotNil(t, tt.expectedLog.Formatter)

			// Reset logrus settings to defaults
			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(logrus.InfoLevel)
			logrus.SetReportCaller(false)
			logrus.SetFormatter(&logrus.TextFormatter{})
		})
	}
}
