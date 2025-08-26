package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"app_aggregator/internal/config"
	"app_aggregator/internal/repository"
	"app_aggregator/internal/router"
	"app_aggregator/internal/services"
	"app_aggregator/migrations"
	"app_aggregator/pkg/closer"
	"app_aggregator/pkg/db"
)

func main() {
	logger := initLogger()
	logger.Info("Starting application")

	closer := closer.NewGracefulCloser()

	logger.Info("Initializing configuration")
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Error("Failed to initialize configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Configuration initialized successfully")

	logger.Info("Running database migrations")
	if err := migrations.Up(cfg); err != nil {
		logger.Error("Failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Database migrations completed")

	logger.Info("Initializing database connection")
	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Database connection established")

	logger.Info("Initializing repositories")
	repo := repository.NewRepository(database)
	organizationRepo := repository.NewOrganizationRepository(repo)
	loanApplicationRepo := repository.NewLoanApplicationsRepository(repo)

	logger.Info("Initializing services")
	organizationService := services.NewOrganizationService(organizationRepo)
	loanApplicationService := services.NewLoanApplicationService(loanApplicationRepo)

	logger.Info("Initializing HTTP server")
	httpServer := router.NewHTTPServer(organizationService, loanApplicationService, logger)

	serverShutdown := make(chan struct{})
	var shutdownOnce sync.Once

	closer.Add(func() error {
		logger.Info("Shutting down HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		select {
		case <-serverShutdown:
			logger.Info("HTTP server already shutdown")
			return nil
		default:
		}

		err := httpServer.Shutdown(ctx)
		shutdownOnce.Do(func() {
			close(serverShutdown)
		})
		return err
	})

	closer.Add(func() error {
		logger.Info("Closing database connections")
		return database.Close()
	})

	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", slog.String("error", err.Error()))
		}
		shutdownOnce.Do(func() {
			close(serverShutdown)
		})
	}()

	logger.Info("Application started successfully", slog.String("port", ":8080"))

	ctx := context.Background()
	closer.Run(ctx, logger)

	logger.Info("Application shutdown completed")
}

func initLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
