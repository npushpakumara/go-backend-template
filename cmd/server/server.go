package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	middlewares "github.com/npushpakumara/go-backend-template/api/middlwares"
	awsclient "github.com/npushpakumara/go-backend-template/internal/aws_client"
	"github.com/npushpakumara/go-backend-template/internal/features/auth"
	"github.com/npushpakumara/go-backend-template/internal/features/email"

	"github.com/gin-gonic/gin"
	"github.com/npushpakumara/go-backend-template/internal/config"
	"github.com/npushpakumara/go-backend-template/internal/features/user"
	"github.com/npushpakumara/go-backend-template/internal/postgres"
	"github.com/npushpakumara/go-backend-template/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Run initializes and starts the application.
// It loads configuration, sets up logging, creates the application container,
// and provides necessary dependencies and services to the application.
func Run() {
	// Load application configuration.
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set up logging with the configuration loaded.
	logging.SetConfig(&logging.Config{
		Encoding:    conf.Logging.Encoding,
		Level:       zapcore.Level(conf.Logging.Level),
		Development: !conf.Server.Production,
	})

	// Ensure that the logger is synced and flushes any pending logs before the application exits.
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logging.DefaultLogger())

	// Create a new application container with various components and configurations.
	app := fx.New(
		// Supply configuration values to the container.
		fx.Supply(conf),
		fx.Supply(conf.AWS.Region),
		fx.Supply(logging.DefaultLogger().Desugar()),
		// Configure the logger for the container.
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		// Set a timeout for graceful shutdown of the application.
		fx.StopTimeout(conf.Server.GracefulShutdown+time.Second),
		// Provide dependencies needed by the application.
		fx.Provide(
			awsclient.NewAWSClient,
			postgres.NewDatabase,
			postgres.NewTransactionManager,
			email.NewSESEmailService,

			// User dependencies
			user.NewUserRepository,
			user.NewUserService,
			user.NewUserHandler,

			// Auth dependencies
			auth.NewAuthService,
			auth.NewAuthHandler,

			middlewares.NewAuthMiddleware,
			newServer,
		),
		// Invoke functions to set up routes and start the application.
		fx.Invoke(
			auth.NewOAuthProviders,
			user.UserRouter,
			auth.AuthRouter,
			func(r *gin.Engine) {},
		),
	)
	// Run the application container.
	app.Run()

}

// newServer creates and configures a new HTTP server using Gin.
// It also sets up lifecycle hooks for starting and stopping the server.
func newServer(lc fx.Lifecycle, cfg *config.Config) *gin.Engine {
	g := gin.New()
	g.Use(gin.Recovery())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      g,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Append hooks to the lifecycle for starting and stopping the server.
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start the server :%d", cfg.Server.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logging.DefaultLogger().Errorw("failed to close http server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Info("Stopped the server")
			return srv.Shutdown(ctx)
		},
	})
	return g
}
