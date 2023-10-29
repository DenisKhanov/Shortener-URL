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
	fmt.Printf("Server started:\nServer addres %s\nBase URL %s\nFile path %s\n", cfg.EnvServAdr, cfg.EnvBaseURL, cfg.EnvStoragePath)
	level, err := logrus.ParseLevel(cfg.EnvLogLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(level)
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

	var myRepository services.Repository
	if cfg.EnvStoragePath == "" {
		myRepository = repositoryes.NewURLInMemoryRepo()
	} else {
		//projectRoot, err := os.Getwd()
		//if err != nil {
		//	log.Fatal(err)
		//}
		////Объединение корневого каталога проекта с подкаталогом tmp и именем файла
		//filePath := filepath.Join(projectRoot, cfg.EnvStoragePath)
		myRepository = repositoryes.NewURLInFileRepo(cfg.EnvStoragePath)
	}

	myService := services.NewServices(myRepository, services.Services{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myService)

	router := mux.NewRouter()
	compressRouter := myHandler.MiddlewareCompress(router)
	loggerRouter := myHandler.MiddlewareLogging(compressRouter)

	router.HandleFunc("/", myHandler.PostURL)
	router.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")
	router.HandleFunc("/api/shorten", myHandler.JSONURL).Methods("POST")
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
