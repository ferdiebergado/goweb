package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/env"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/ferdiebergado/goweb/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfgFile := flag.String("config", "config.json", "Config file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx, *cfgFile); err != nil {
		slog.Error("fatal error", "reason", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfgFile string) error {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer func() {
		stop()
		slog.Info("Signal context cleanup complete.")
	}()

	config, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("load environment: %w", err)
	}

	setLogger(os.Stdout, &config.App)

	db, err := openDB(&config.Db)
	if err != nil {
		return err
	}

	defer db.Close()

	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Db.PingTimeout)*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	db.SetMaxOpenConns(30)

	slog.Info("Connected to the database", "db", config.Db.DB)

	router := goexpress.New()
	setupRoutes(router, db)

	server := &http.Server{ // #nosec G112 -- timeouts will be handled by reverse proxy
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a separate goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Server started", "address", server.Addr, "env", config.App.Env)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	// Wait for a shutdown signal or server error
	select {
	case <-signalCtx.Done(): // Received termination signal (CTRL+C)
		slog.Info("Shutdown signal received.")
	case err := <-serverErr: // Server crashed
		return fmt.Errorf("server error: %w", err)
	}

	// Graceful shutdown with timeout
	slog.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("Server gracefully shut down.")
	return nil
}

func loadEnv() (string, error) {
	const dev = "development"
	var envFile string
	appEnv := env.Get("ENV", dev)

	switch appEnv {
	case "production":
		return appEnv, nil
	case dev:
		envFile = ".env"
	case "testing":
		envFile = ".env.testing"
	default:
		return "", fmt.Errorf("unrecognized environment: %s", appEnv)
	}

	if err := env.Load(envFile); err != nil {
		return "", fmt.Errorf("cannot load env file: %s", envFile)
	}

	return appEnv, nil
}

func setLogger(out io.Writer, cfg *config.AppConfig) {
	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		if cfg.IsDebug {
			logLevel.Set(slog.LevelDebug)
		}

		opts.AddSource = true
		handler = slog.NewTextHandler(out, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func openDB(cfg *config.DBConfig) (*sql.DB, error) {
	slog.Info("Connecting to the database")
	const dbStr = "postgres://%s:%s@localhost:5432/%s?sslmode=disable"
	dsn := fmt.Sprintf(dbStr, cfg.User, cfg.Pass, cfg.DB)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("database initialization: %w", err)
	}

	return db, nil
}

func setupRoutes(r *goexpress.Router, db *sql.DB) {
	r.Use(goexpress.RecoverFromPanic)
	r.Use(goexpress.LogRequest)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	baseHandler := handler.NewBaseHandler(service)
	mountRoutes(r, baseHandler)
}
