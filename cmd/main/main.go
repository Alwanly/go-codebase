package main

import (
	"context"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/Alwanly/go-codebase/config"
	"github.com/Alwanly/go-codebase/pkg/authentication"
	"github.com/Alwanly/go-codebase/pkg/database"
	"github.com/Alwanly/go-codebase/pkg/logger"
	"github.com/Alwanly/go-codebase/pkg/middleware"
	"github.com/Alwanly/go-codebase/pkg/redis"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// @title Codebase API Example documentation
// @version 1.0
// @description This is a sample server.

// @host localhost:9000
// @BasePath /

// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// load config
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup dependencies
	globalLogger := logger.NewLogger(cfg.ServiceName, cfg.LogLevel,
		logger.WithPrettyPrint(),
	)
	l := logger.WithID(globalLogger, "server", "main")
	l.Info("Starting application",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
	)

	// Setup database
	dbConfig := database.DBServiceOpts{
		Debug:                      cfg.Debug,
		Logger:                     globalLogger,
		PostgresURI:                &cfg.PostgresURI,
		PostgresMaxOpenConnections: cfg.PostgresMaxOpenConnections,
		PostgresMaxIdleConnections: cfg.PostgresMaxIdleConnections,
	}

	db, err := database.NewPostgres(&dbConfig)
	if err != nil {
		l.Error("Failed to initialize database", zap.Error(err))
		os.Exit(1)
	}

	// Setup redis
	redisConfig := redis.Opts{
		Logger:   globalLogger,
		RedisURI: &cfg.RedisURI,
	}
	redisClient, err := redis.NewRedis(&redisConfig)
	if err != nil {
		l.Error("Failed to initialize Redis", zap.Error(err))
		os.Exit(1)
	}

	// Setup middleware
	jwtConfig := middleware.SetJwtAuth(&authentication.JWTConfig{
		PrivateKey:     cfg.PrivateKey,
		PublicKey:      cfg.PublicKey,
		Audience:       cfg.JwtAudience,
		Issuer:         cfg.JwtIssuer,
		ExpirationTime: cfg.JwtExpirationTime,
		RefreshTime:    cfg.JwtRefreshTime,
	})
	basicAuthConfig := middleware.SetBasicAuth(&authentication.BasicAuthTConfig{
		Username: cfg.BasicAuthUsername,
		Password: cfg.BasicAuthPassword,
	})

	authMiddleware := middleware.NewAuthMiddleware(jwtConfig, basicAuthConfig)
	if authMiddleware == nil {
		l.Error("Failed to create auth middleware")
		os.Exit(1)
	}

	// Create app
	app := Bootstrap(&AppDeps{
		Config: &cfg,
		Logger: globalLogger,
		DB:     db,
		Redis:  redisClient,
		Auth:   authMiddleware,
	})

	// Register health check

	//--------------------- Bootstrap Application ---------------------

	ctx, cancel := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)

	// run http server
	g.Go(func() error {
		l.Info("Starting server...", zap.Int("port", cfg.Port))
		return app.Fiber.Listen(fmt.Sprintf(":%d", cfg.Port))
	})

	// graceful shutdown
	g.Go(func() error {
		<-gCtx.Done()
		l.Info("Gracefully shutting down...")

		l.Info("Server gracefully shutdown")
		if err := app.Fiber.Shutdown(); err != nil {
			l.Error("Cannot shutdown server", zap.Error(err))
			return err
		}

		l.Info("Closing database connection")
		if err := app.DB.Close(); err != nil {
			l.Error("Cannot close database connection", zap.Error(err))
			return err
		}

		return nil
	})

	// listen for interrupt signal
	go func() {
		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		l.Info("Listening for OS signal...")
		<-c

		// cancel context
		l.Info("Received OS signal, canceling context...")
		cancel()
	}()

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
