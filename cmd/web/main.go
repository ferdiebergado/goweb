package main

import (
	"context"
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
	"github.com/ferdiebergado/goweb/internal/infra/db"
	"github.com/ferdiebergado/goweb/internal/pkg/security"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var validate *validator.Validate

func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		stop()
		slog.Info("Signal context cleanup complete.")
	}()

	if err := run(signalCtx); err != nil {
		slog.Error("fatal error", "reason", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	const (
		envVar  = "ENV"
		envDev  = "development"
		envProd = "production"
		cfgFile = "config.json"
		fmtAddr = ":%d"
	)
	appEnv := env.Get(envVar, envDev)
	if appEnv != envProd {
		if err := loadEnv(appEnv); err != nil {
			return fmt.Errorf("load env: %w", err)
		}
	}

	setLogger(os.Stdout, appEnv)

	cf := flag.String("cfg", cfgFile, "Config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cf)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := db.Connect(ctx, &cfg.Db)
	if err != nil {
		return err
	}
	defer db.Close()

	router := goexpress.New()
	configureValidator()
	tmpl, err := handler.NewTemplate(cfg.Template)
	if err != nil {
		return err
	}
	hasher := &security.Argon2Hasher{}

	deps := &handler.AppDependencies{
		Config:    cfg,
		DB:        db,
		Router:    router,
		Validator: validate,
		Template:  tmpl,
		Hasher:    hasher,
	}

	app := handler.NewApp(deps)
	app.SetupRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(fmtAddr, cfg.Server.Port),
		Handler:      app.Router(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Run server in a separate goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Server started", "address", server.Addr, "env", cfg.App.Env, slog.Bool("debug", cfg.App.IsDebug))
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

func setLogger(out io.Writer, appEnv string) {
	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	if appEnv == "production" {
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

func configureValidator() {
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
