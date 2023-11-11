package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/config"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/logcfg"
	"github.com/DenisKhanov/shorterURL/internal/app/repositoryes"
	"github.com/DenisKhanov/shorterURL/internal/app/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		dbPool       *pgxpool.Pool
		err          error
		cfg          *config.ENVConfig
		myRepository services.Repository
	)
	cfg = config.NewConfig()
	if cfg.EnvDataBase != "" {
		dbPool, err = pgxpool.New(context.Background(), cfg.EnvDataBase)
		if err != nil {
			logrus.Error("Don't connect to DB: ", err)
			os.Exit(1)
		}
		defer dbPool.Close()
		myRepository = repositoryes.NewURLInDBRepo(dbPool)
	} else {
		myRepository = repositoryes.NewURLInMemoryRepo(cfg.EnvStoragePath)
	}

	logcfg.RunLoggerConfig(cfg.EnvLogLevel)
	logrus.Infof("Server started:\nServer addres %s\nBase URL %s\nFile path %s\nDBConfig %s\n", cfg.EnvServAdr, cfg.EnvBaseURL, cfg.EnvStoragePath, cfg.EnvDataBase)

	myService := services.NewServices(myRepository, services.Services{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myService, dbPool)

	router := mux.NewRouter()
	compressRouter := myHandler.MiddlewareCompress(router)
	loggerRouter := myHandler.MiddlewareLogging(compressRouter)

	router.HandleFunc("/", myHandler.PostURL)
	router.HandleFunc("/ping", myHandler.PingDB)
	router.HandleFunc("/{id}", myHandler.GetURL).Methods("GET")
	router.HandleFunc("/api/shorten", myHandler.JSONURL).Methods("POST")
	router.HandleFunc("/api/shorten/batch", myHandler.BatchSave).Methods("POST")

	server := &http.Server{Addr: cfg.EnvServAdr, Handler: loggerRouter}

	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Error(err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP server Shutdown: %v\n", err)
	}
	if inMemoryRepo, ok := myRepository.(services.URLInMemoryRepository); ok {
		err = inMemoryRepo.SaveBatchToFile()
		if err != nil {
			logrus.Error(err)
		}
	}
	logrus.Info("Server exited")
}
