package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("Fatal error.", "reason", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	return nil
}
