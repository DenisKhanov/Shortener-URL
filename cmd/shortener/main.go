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

//Задание по треку «Сервис сокращения URL»
//Добавьте в сервис функциональность подключения к базе данных. В качестве СУБД используйте PostgreSQL не ниже 10 версии.
//Добавьте в сервис хендлер GET /ping, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
//Строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN или флага командной строки -d.
//Для работы с БД используйте один из следующих пакетов:
//database/sql,
//github.com/jackc/pgx,
//github.com/lib/pq,
//github.com/jmoiron/sqlx.

var cfg *config.ENVConfig

func main() {
	cfg = config.NewConfig()
	fmt.Printf("Server started:\nServer addres %s\nBase URL %s\nFile path %s\nDBConfig %s\n", cfg.EnvServAdr, cfg.EnvBaseURL, cfg.EnvStoragePath, cfg.EnvDataBase)

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

	myRepository := repositoryes.NewURLInMemoryRepo(cfg.EnvStoragePath)
	myService := services.NewServices(myRepository, services.Services{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myService, cfg.EnvDataBase)

	router := mux.NewRouter()
	compressRouter := myHandler.MiddlewareCompress(router)
	loggerRouter := myHandler.MiddlewareLogging(compressRouter)

	router.HandleFunc("/", myHandler.PostURL)
	router.HandleFunc("/ping", myHandler.ConfigDB)
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

	if err = server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP server Shutdown: %v\n", err)
	}
	err = myRepository.SaveBatchToFile()
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("Server exited")
}
