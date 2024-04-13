package app

import (
	"context"
	myGRPC "github.com/DenisKhanov/shorterURL/internal/api/grpc/interceptors"
	"github.com/DenisKhanov/shorterURL/internal/api/http/middleware"
	"github.com/DenisKhanov/shorterURL/internal/config"
	"github.com/DenisKhanov/shorterURL/internal/logcfg"
	"github.com/DenisKhanov/shorterURL/internal/models"
	"github.com/DenisKhanov/shorterURL/internal/services/url"
	"github.com/DenisKhanov/shorterURL/internal/tls"
	proto "github.com/DenisKhanov/shorterURL/pkg/shortener_v1"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	realip "github.com/thanhhh/gin-gonic-realip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// App represents the application structure responsible for initializing dependencies
// and running the http_shortener serverHTTP.
type App struct {
	serviceProvider *serviceProvider  // The service provider for dependency injection
	config          *config.ENVConfig // The configuration object for the application
	dbPool          *pgxpool.Pool     // The connection pool to the database
	trustedSubnets  []*net.IPNet      //The collection trusted subnet
	serverHTTP      *http.Server      // The serverHTTP instance
	serverGRPC      *grpc.Server      //The serverGRPC instance
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

// Run starts the application and runs the http_shortener serverHTTP.
func (a *App) Run() {
	a.runShortenerServers()
}

// initDeps initializes all dependencies required by the application.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initTrustedSubnets,
		a.initDBConnection,
		a.initServiceProvider,
		a.initShortenerHTTPServer,
		a.initShortenerGRPCServer,
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

// parseSubnets parses a string containing a list of CIDR subnets and returns them as a []*net.IPNet objects.
func (a *App) initTrustedSubnets(_ context.Context) error {
	var subnets []*net.IPNet
	if a.config.EnvSubnet != "" {
		subStr := strings.Split(a.config.EnvSubnet, ",")
		for _, subnetStr := range subStr {
			_, subnetIPNet, err := net.ParseCIDR(subnetStr)
			if err != nil {
				logrus.WithError(err).Error("error parsing string CIDR")
				return err
			}
			subnets = append(subnets, subnetIPNet)
		}
	}
	a.trustedSubnets = subnets
	return nil
}

// initDBConnection initializes the connection to the database.
func (a *App) initDBConnection(ctx context.Context) error {
	if a.config.EnvDataBase != "" {
		confPool, err := pgxpool.ParseConfig(a.config.EnvDataBase)
		if err != nil {
			logrus.WithError(err).Error("Error parsing config")
			return err
		}
		confPool.MaxConns = 50
		confPool.MinConns = 10
		a.dbPool, err = pgxpool.NewWithConfig(ctx, confPool)
		if err != nil {
			logrus.WithError(err).Error("Don't connect to DB")
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

// initShortenerHTTPServer initializes the http_shortener serverHTTP with middleware and routes.
func (a *App) initShortenerHTTPServer(_ context.Context) error {
	myHandler := a.serviceProvider.ShortenerHandler(a.dbPool, a.config.EnvBaseURL, a.config.EnvStoragePath, a.config.EnvSubnet)

	// Установка переменной окружения для включения режима разработки
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	// Use the pprof middleware
	pprof.Register(router)
	//Public middleware routers group
	publicRoutes := router.Group("/")
	publicRoutes.Use(middleware.AuthPublic())
	publicRoutes.Use(middleware.LogrusLog())
	publicRoutes.Use(middleware.GZIPCompress())

	publicRoutes.POST("/", myHandler.GetShortURL)
	publicRoutes.GET("/ping", myHandler.GetStorageStatus)
	publicRoutes.GET("/:id", myHandler.GetOriginalURL)
	publicRoutes.POST("/api/shorten", myHandler.GetJSONShortURL)
	publicRoutes.POST("/api/shorten/batch", myHandler.GetBatchShortURL)

	//Private middleware routers group
	privateRoutes := router.Group("/")
	privateRoutes.Use(middleware.AuthPrivate())
	privateRoutes.Use(middleware.LogrusLog())
	privateRoutes.Use(middleware.GZIPCompress())

	privateRoutes.GET("/api/user/urls", myHandler.GetUserURLS)
	privateRoutes.DELETE("/api/user/urls", myHandler.DelUserURLs)

	//Only trusted subnet middleware
	trustSubnetRouter := router.Group("/")
	trustSubnetRouter.Use(realip.RealIP())
	trustSubnetRouter.Use(middleware.TrustedSubnet(a.trustedSubnets))

	trustSubnetRouter.GET("/api/internal/stats", myHandler.GetServiceStats)

	a.serverHTTP = &http.Server{
		Addr:    a.config.EnvServAdr,
		Handler: router,
	}

	return nil
}

// initShortenerGRPCServer initializes the  serverGRPC with interceptors.
func (a *App) initShortenerGRPCServer(_ context.Context) error {
	newGRPC := a.serviceProvider.ShortenerGRPC()

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(myGRPC.UnaryLoggerInterceptor,
		myGRPC.UnaryTrustedSubnetsInterceptor(a.trustedSubnets),
		myGRPC.UnaryPrivateAuthInterceptor, myGRPC.UnaryPublicAuthInterceptor),
	)
	reflection.Register(s)
	a.serverGRPC = s
	// регистрируем сервис
	proto.RegisterShortenerV1Server(s, newGRPC)
	return nil

}

// runShortenerServers starts the gRPC + HTTP servers with graceful shutdown.
func (a *App) runShortenerServers() {
	logcfg.RunLoggerConfig(a.config.EnvLogLevel)

	//run gRPC server
	go func() {
		listen, err := net.Listen("tcp", a.config.EnvGRPC)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Infof("Starting server gRPC on: %s", a.config.EnvGRPC)
		if err = a.serverGRPC.Serve(listen); err != nil {
			logrus.WithError(err).Error("The server gRPC  failed to start")
		}
	}()

	// run HTTP server
	go func() {
		if a.config.EnvTLS != "" {
			logrus.Infof("Starting server HTTP with TLS on: %s", a.config.EnvServAdr)
			_, err := tls.NewTLS()
			if err != nil {
				logrus.WithError(err).Error("Error generate TLS")
			}
			if err = a.serverHTTP.ListenAndServeTLS(models.CertPEM, models.PrivateKeyPEM); err != nil {
				logrus.Error(err)
			}
		} else {
			logrus.Infof("Starting serverHTTP on: %s", a.config.EnvServAdr)
			if err := a.serverHTTP.ListenAndServe(); err != nil {
				logrus.Error(err)
			}
		}
	}()

	// Shutdown signal with grace period of 5 seconds
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan
	logrus.Infof("Shutting down HTTP & gRPC servers with signal : %v...", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		a.serverGRPC.GracefulStop()
		logrus.Infof("gRPC servers closed")
		wg.Done()
	}()

	go func() {
		if err := a.serverHTTP.Shutdown(shutdownCtx); err != nil {
			logrus.WithError(err).Error("HTTP server shutdown error")
		}
		wg.Done()
	}()

	//TODO избавиться от приведения типов

	//If the input shutdown signal, batch URLs saving to file
	if a.dbPool == nil {
		err := a.serviceProvider.shortenerRepository.(url.InMemoryRepository).SaveBatchToFile()
		if err != nil {
			logrus.WithError(err).Error("Error save memory in file")
		}
	} else {
		a.dbPool.Close()
	}
	wg.Wait()
	logrus.Info("Server exited")
}
