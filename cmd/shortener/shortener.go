package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DenisKhanov/shorterURL/internal/app/config"
	"github.com/DenisKhanov/shorterURL/internal/app/handlers"
	"github.com/DenisKhanov/shorterURL/internal/app/logcfg"
	"github.com/DenisKhanov/shorterURL/internal/app/repositories"
	"github.com/DenisKhanov/shorterURL/internal/app/services"
	"github.com/gin-contrib/pprof" // подключаем пакет pprof gin
	"github.com/gin-gonic/gin"
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
			logrus.Errorf("error parsing config: %v", err)
		}
		confPool.MaxConns = 50
		confPool.MinConns = 10
		dbPool, err = pgxpool.NewWithConfig(context.Background(), confPool)
		if err != nil {
			logrus.Error("Don't connect to DB: ", err)
			logrus.Fatal(err)
		}

		defer dbPool.Close()
		myRepository = repositories.NewURLInDBRepo(dbPool)
	} else {
		myRepository = repositories.NewURLInMemoryRepo(cfg.EnvStoragePath)
		repositoryReciver = true
	}

	logcfg.RunLoggerConfig(cfg.EnvLogLevel)
	logrus.Infof("Server started:\nServer addres %s\nBase URL %s\nFile path %s\nDBConfig %s\n", cfg.EnvServAdr, cfg.EnvBaseURL, cfg.EnvStoragePath, cfg.EnvDataBase)
	myShorURLService := services.NewShortURLServices(myRepository, services.ShortURLServices{}, cfg.EnvBaseURL)
	myHandler := handlers.NewHandlers(myShorURLService, dbPool)

	// Установка переменной окружения для отключения режима разработки
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	// Use the pprof middleware
	pprof.Register(router)
	//Public middleware routers group
	publicRoutes := router.Group("/")
	publicRoutes.Use(myHandler.MiddlewareAuthPublic())
	publicRoutes.Use(myHandler.MiddlewareLogging())
	publicRoutes.Use(myHandler.MiddlewareCompress())

	publicRoutes.POST("/", myHandler.GetShortURL)
	publicRoutes.GET("/ping", myHandler.PingDB)
	publicRoutes.GET("/:id", myHandler.GetOriginalURL)
	publicRoutes.POST("/api/shorten", myHandler.GetJSONShortURL)
	publicRoutes.POST("/api/shorten/batch", myHandler.GetBatchShortURL)
	//Private middleware routers group
	privateRoutes := router.Group("/")
	privateRoutes.Use(myHandler.MiddlewareAuthPrivate())
	privateRoutes.Use(myHandler.MiddlewareLogging())
	privateRoutes.Use(myHandler.MiddlewareCompress())

	privateRoutes.GET("/api/user/urls", myHandler.GetUserURLS)
	privateRoutes.DELETE("/api/user/urls", myHandler.DelUserURLS)

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
