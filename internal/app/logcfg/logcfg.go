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

func RunLoggerConfig(EnvLogs string) {

	logLevel, err := logrus.ParseLevel(EnvLogs)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s.%d.%s", filename, f.Line, f.Function)
			return "", filename
		},
	})
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "shortener.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
	})
	logrus.SetOutput(mw)
}
