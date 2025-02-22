package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-service/internal/config"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository/postgres"
	"wallet-service/internal/router"
	"wallet-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	fmt.Println(cfg)

	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.NewDB(cfg.DBConnStr)

	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repository := postgres.NewRepository(db, *logger)
	service := service.NewService(repository, *logger)
	handler := handler.NewHandler(service, *logger)

	app := router.SetupRouter(*handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info(fmt.Sprintf("Server started on port %s", cfg.Port))
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown: ", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")

}
