package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/abibby/eztvrss/app"
	"github.com/abibby/salusa/clog"
)

func main() {
	ctx := clog.Init(context.Background(), slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})))
	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		clog.Use(ctx).Error("error bootstrapping", "error", err)
		os.Exit(1)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		clog.Use(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
