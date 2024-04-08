package app

import (
	"context"
	"errors"
	"github.com/DenisKhanov/shorterURL/internal/config"
	"github.com/DenisKhanov/shorterURL/internal/https"
	"github.com/DenisKhanov/shorterURL/internal/logcfg"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/DenisKhanov/shorterURL/internal/services"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	realip "github.com/thanhhh/gin-gonic-realip"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App represents the application structure responsible for initializing dependencies
// and running the shortener server.
type App struct {
	serviceProvider *serviceProvider  // The service provider for dependency injection
	config          *config.ENVConfig // The configuration object for the application
	dbPool          *pgxpool.Pool     // The connection pool to the database
	server          *http.Server      // The HTTP server instance
}

// NewApp creates a new instance of the application.
func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Run starts the application and runs the shortener server.
func (a *App) Run() error {
	return a.runShortenerServer()
}

// initDeps initializes all dependencies required by the application.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initDBConnection,
		a.initServiceProvider,
		a.initShortenerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig initializes the application configuration.
func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	a.config = cfg
	config.PrintProjectInfo()
	return nil
}

// initDBConnection initializes the connection to the database.
func (a *App) initDBConnection(ctx context.Context) error {
	if a.config.EnvDataBase != "" {
		confPool, err := pgxpool.ParseConfig(a.config.EnvDataBase)
		if err != nil {
			logrus.Errorf("error parsing config: %v", err)
			return err
		}
		confPool.MaxConns = 50
		confPool.MinConns = 10
		a.dbPool, err = pgxpool.NewWithConfig(ctx, confPool)
		if err != nil {
			logrus.Error("Don't connect to DB: ", err)
			return err
		}
	}
	return nil
}

// initServiceProvider initializes the service provider for dependency injection.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initShortenerServer initializes the shortener server with middleware and routes.
func (a *App) initShortenerServer(_ context.Context) error {
	logcfg.RunLoggerConfig(a.config.EnvLogLevel)

	myHandler, err := a.serviceProvider.UserHandler(a.dbPool, a.config.EnvBaseURL, a.config.EnvStoragePath, a.config.EnvSubnet)
	if err != nil {
		logrus.Error(err)
	}

	// Установка переменной окружения для включения режима разработки
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

	//Only trusted subnet middleware
	trustSubnetRouter := router.Group("/api/internal")
	trustSubnetRouter.Use(realip.RealIP())
	trustSubnetRouter.Use(myHandler.MiddlewareTrustedSubnets())

	trustSubnetRouter.GET("/api/internal/stats", myHandler.ServiceStats)

	a.server = &http.Server{
		Addr:    a.config.EnvServAdr,
		Handler: router,
	}

	return nil
}

// runShortenerServer starts the shortener server and handles graceful shutdown.
func (a *App) runShortenerServer() error {
	defer a.dbPool.Close()
	go func() {
		if a.config.EnvHTTPS != "" {
			logrus.Info("Starting server with TLS on: ", a.config.EnvServAdr)
			_, err := https.NewHTTPS()
			if err != nil {
				logrus.Error(err)
			}
			if err = a.server.ListenAndServeTLS(models.CertPEM, models.PrivateKeyPEM); !errors.Is(err, http.ErrServerClosed) {
				logrus.Fatal(err)
			}
		} else {
			logrus.Info("Starting server on: ", a.config.EnvServAdr)
			if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				logrus.Error(err)
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	logrus.Infof("Shutting down server with signal : %v...", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		logrus.Errorf("HTTP server Shutdown error: %v\n", err)
	}
	//If the server shutting down, save batch to file
	if a.dbPool == nil {
		err := a.serviceProvider.userRepository.(services.URLInMemoryRepository).SaveBatchToFile()
		if err != nil {
			logrus.Error(err)
		}
	}

	logrus.Info("Server exited")
	return nil
}
