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
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/env"
	"github.com/ferdiebergado/goweb/internal/config"
	"github.com/ferdiebergado/goweb/internal/handler"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var validate *validator.Validate

func main() {
	setLogger(os.Stdout)

	cfgFile := flag.String("cfg", "config.json", "Config file")
	flag.Parse()

	if appEnv := env.Get("ENV", "development"); appEnv != "production" {
		if err := loadEnv(appEnv); err != nil {
			logFatal(fmt.Errorf("load env: %w", err))
		}
	}

	cfg, err := config.LoadConfig(*cfgFile)
	if err != nil {
		logFatal(fmt.Errorf("load config: %w", err))
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		stop()
		slog.Info("Signal context cleanup complete.")
	}()

	if err := run(signalCtx, cfg); err != nil {
		logFatal(err)
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	db, err := openDB(ctx, &cfg.Db)
	if err != nil {
		return err
	}
	defer db.Close()

	router := goexpress.New()
	setValidate()
	app := handler.NewApp(cfg, db, router, validate)
	app.SetupRoutes()

	server := &http.Server{ // #nosec G112 -- timeouts will be handled by reverse proxy
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: app.Router(),
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
	case <-ctx.Done(): // Received termination signal (CTRL+C)
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

func loadEnv(appEnv string) error {
	const (
		envDev  = ".env"
		envTest = ".env.testing"
	)
	var envFile string

	switch appEnv {
	case "development":
		envFile = envDev
	case "testing":
		envFile = envTest
	default:
		return fmt.Errorf("unrecognized environment: %s", appEnv)
	}

	if err := env.Load(envFile); err != nil {
		return fmt.Errorf("cannot load env file %s, %w", envFile, err)
	}

	return nil
}

func setLogger(out io.Writer) {
	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	e := os.Getenv("ENV")

	if e == "production" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		if env.GetBool("DEBUG", false) {
			logLevel.Set(slog.LevelDebug)
		}

		opts.AddSource = true
		handler = slog.NewTextHandler(out, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func openDB(ctx context.Context, cfg *config.DBConfig) (*sql.DB, error) {
	const dbStr = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
	slog.Info("Connecting to the database")
	dsn := fmt.Sprintf(dbStr, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
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

func setValidate() {
	validate = validator.New()

	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func logFatal(err error) {
	slog.Error("fatal error", "reason", err)
	os.Exit(1)
}
