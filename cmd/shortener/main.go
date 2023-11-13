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
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		dbPool            *pgxpool.Pool
		err               error
		cfg               *config.ENVConfig
		myRepository      services.Repository
		repositoryReciver bool
	)

	cfg = config.NewConfig()
	if cfg.EnvDataBase != "" {
		confPool, err := pgxpool.ParseConfig(cfg.EnvDataBase)
		if err != nil {
			log.Fatalf("error parsing config: %v", err)
		}
		confPool.MaxConns = 50
		confPool.MinConns = 10
		dbPool, err = pgxpool.NewWithConfig(context.Background(), confPool)
		if err != nil {
			logrus.Error("Don't connect to DB: ", err)
			os.Exit(1)
		}

		defer dbPool.Close()
		myRepository = repositoryes.NewURLInDBRepo(dbPool)
	} else {
		myRepository = repositoryes.NewURLInMemoryRepo(cfg.EnvStoragePath)
		repositoryReciver = true
	}

	logcfg.RunLoggerConfig(cfg.EnvLogLevel)
	logrus.Infof("Server started:\nServer addres %s\nBase URL %s\nFile path %s\nDBConfig %s\n", cfg.EnvServAdr, cfg.EnvBaseURL, cfg.EnvStoragePath, cfg.EnvDataBase)
	myShorURLService := services.NewShortURLServices(myRepository, services.ShortURLServices{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myShorURLService, dbPool)

	router := gin.Default()
	router.Use(myHandler.MiddlewareLogging())
	router.Use(gzip.Gzip(gzip.BestSpeed))

	router.POST("/", myHandler.GetShortURL)
	router.GET("/ping", myHandler.PingDB)
	router.GET("/{id}", myHandler.GetOriginalURL)
	router.POST("/api/shorten", myHandler.GetJSONShortURL)
	router.POST("/api/shorten/batch", myHandler.GetBatchJSONShortURL)

	server := &http.Server{Addr: cfg.EnvServAdr, Handler: router}

	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	go func() {
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
	//If the server shutting down, save batch to file
	if repositoryReciver {
		err = myRepository.(services.URLInMemoryRepository).SaveBatchToFile()
		if err != nil {
			logrus.Error(err)
		}
	}

	logrus.Info("Server exited")
}
