package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/config"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/repositoryes"
	"github.com/DenisKhanov/shorterURL/internal/app/services"
	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
)

var cfg *config.ENVConfig

func init() {
	cfg = config.NewConfig()
}
func main() {
	fmt.Println("Server Address:", cfg.EnvServAdr)
	fmt.Println("Base URL:", cfg.EnvBaseURL)
	//fmt.Println("Log level:", cfg.EnvLogLevel)

	//level, err := logrus.ParseLevel(cfg.EnvLogLevel)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s:%d", filename, f.Line)
			return "", filename
		},
	})
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, //mb
		MaxBackups: 3,
		MaxAge:     30, //day
	})

	myRepository := repositoryes.NewRepository(make(map[string]string), make(map[string]string))
	myService := services.NewServices(myRepository, services.Services{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myService)

	r := mux.NewRouter()
	compressRouter := myHandler.MiddlewareCompress(r)
	loggerRouter := myHandler.MiddlewareLogging(compressRouter)

	r.HandleFunc("/", myHandler.PostURL)
	r.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")
	r.HandleFunc("/api/shorten", myHandler.JSONURL).Methods("POST")
	server := &http.Server{Addr: cfg.EnvServAdr, Handler: loggerRouter}

	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP server Shutdown: %v\n", err)
	}
	fmt.Println("Server exited")

}
