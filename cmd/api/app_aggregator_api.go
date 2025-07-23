package main

import (
	"app_aggregator/internal/config"
	"app_aggregator/internal/repository"
	"app_aggregator/internal/router"
	"app_aggregator/migrations"
	"app_aggregator/pkg/db"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func main() {
	gin.SetMode(gin.DebugMode)
	logger := loggerInit()

	logger.Info("Configuration initialization has started")
	cfg, err := config.InitConfig()
	if err != nil {
		logger.Error("configuration initialization has failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Configuration initialization has finished")

	logger.Info("Database migration has started")
	err = migrations.Up(cfg)
	if err != nil {
		logger.Error("migration has failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Database migration has finished")

	logger.Info("Database initialization has started")
	pgDb, err := db.InitDB(cfg)
	if err != nil {
		logger.Error("database initialization has failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Database initialization has finished")

	repo := repository.NewRepository(pgDb)

	routerBuilder := router.NewBuilder(repo)
	routerBuilder.OrganizationRouter()

	logger.Info("Server initialization has started")
	if err := routerBuilder.GetEngine().Run(":8080"); err != nil {
		logger.Error("server initialization has failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Server initialization has started at :8080")

}

func loggerInit() *slog.Logger {
	loggerHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)
	return logger
}
