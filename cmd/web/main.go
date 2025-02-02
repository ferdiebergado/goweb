package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("Fatal error.", "reason", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	fmt.Println("App started. Press Ctrl-C to exit.")
	<-signalCtx.Done()

	return nil
}
