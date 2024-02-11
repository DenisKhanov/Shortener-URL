// Package logcfg provides configuration for the application logger.
// It includes functionality to set the log level, report caller information, and configure log file rotation.
package logcfg

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

// RunLoggerConfig configures the application logger based on the provided log level.
// It sets the log level, enables reporting of caller information, and configures log file rotation.
//
// Parameters:
//   - EnvLogs: The log level to set, provided as a string.
//     Valid log levels are "panic", "fatal", "error", "warn", "info", and "debug".
func RunLoggerConfig(EnvLogs string) {
	// Parse log level from the environment variable
	logLevel, err := logrus.ParseLevel(EnvLogs)
	if err != nil {
		logrus.Fatal(err)
	}

	// Set log level and enable reporting of caller information
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	// Customize the log formatter to include line numbers in filenames
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s.%d.%s", filename, f.Line, f.Function)
			return "", filename
		},
	})

	// Configure log file rotation using lumberjack
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "shortener.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
	})
	logrus.SetOutput(mw)
}
