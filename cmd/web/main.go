package main

import (
	"context"
	"database/sql"
	"errors"
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
	ctx := context.WithoutCancel(context.Background())

	if err := run(ctx); err != nil {
		logFatal("Fatal error.", err)
	}
}

func run(ctx context.Context) error {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	appEnv, err := loadEnv()
	if err != nil {
		logFatal("failed to load environment", err)
	}

	setLogger(os.Stdout)

	db := openDB()
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		logFatal("Cannot connect to the database", err)
	}

	db.SetMaxOpenConns(30)

	slog.Info("Connected to the database", "db", os.Getenv("POSTGRES_DB"))

	router := goexpress.New()
	setupRoutes(router, db)

	server := &http.Server{ // #nosec G112 -- timeouts will be handled by reverse proxy
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		slog.Info("Server started", "address", server.Addr, "env", appEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "reason", err)
		}
	}()

	// Wait for shutdown signal
	<-signalCtx.Done()

	slog.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
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

func setLogger(out io.Writer) {
	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		if isDebug := env.GetBool("DEBUG", false); isDebug {
			logLevel.Set(slog.LevelDebug)
		}

		opts.AddSource = true
		handler = slog.NewTextHandler(out, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func openDB() *sql.DB {
	slog.Info("Connecting to the database")
	const dbStr = "postgres://%s:%s@localhost:5432/%s?sslmode=disable"
	dsn := fmt.Sprintf(dbStr, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logFatal("Database initialization failed", err)
	}

	return db
}

func setupRoutes(r *goexpress.Router, db *sql.DB) {
	r.Use(goexpress.RecoverFromPanic)
	r.Use(goexpress.LogRequest)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	baseHandler := handler.NewBaseHandler(service)
	mountRoutes(r, baseHandler)
}

func logFatal(msg string, err error) {
	slog.Error(msg, "reason", err)
	os.Exit(1)
}
