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
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/ferdiebergado/goweb/internal/repository"
	"github.com/ferdiebergado/goweb/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfgFile := flag.String("cfg", "config.json", "Config file")
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

	if err := loadEnv(); err != nil {
		return fmt.Errorf("load env: %w", err)
	}

	cfg, err := loadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	setLogger(os.Stdout, &cfg.App)

	db, err := openDB(ctx, &cfg.Db)
	if err != nil {
		return err
	}
	defer db.Close()

	router := goexpress.New()
	setupRoutes(router, db)

	server := &http.Server{ // #nosec G112 -- timeouts will be handled by reverse proxy
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Run server in a separate goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Server started", "address", server.Addr, "env", cfg.App.Env)
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("Server gracefully shut down.")
	return nil
}

func loadEnv() error {
	const (
		dev     = "development"
		envDev  = ".env"
		envTest = ".env.testing"
	)
	var envFile string
	appEnv := env.Get("ENV", dev)

	switch appEnv {
	case "production":
		return nil
	case dev:
		envFile = envDev
	case "testing":
		envFile = envTest
	default:
		return fmt.Errorf("unrecognized environment: %s", appEnv)
	}

	if err := env.Load(envFile); err != nil {
		return fmt.Errorf("cannot load env file: %s", envFile)
	}

	return nil
}

func setLogger(out io.Writer, cfg *AppConfig) {
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

func openDB(ctx context.Context, cfg *DBConfig) (*sql.DB, error) {
	const dbStr = "postgres://%s:%s@localhost:5432/%s?sslmode=disable"
	slog.Info("Connecting to the database")
	dsn := fmt.Sprintf(dbStr, cfg.User, cfg.Pass, cfg.DB)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("database initialization: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.PingTimeout)*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdle) * time.Second)

	slog.Info("Connected to the database", "db", cfg.DB)
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
