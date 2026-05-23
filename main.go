// main.go
package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/wreckitral/production-backend-go/internal/platform/config"
	"github.com/wreckitral/production-backend-go/internal/platform/db"
	"github.com/wreckitral/production-backend-go/internal/platform/server"
)

//go:embed web/*
var webFS embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := newLogger(cfg.Log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer pool.Close()

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("bind %s: %w", addr, err)
	}

	app := server.New(server.Deps{
		Config:   cfg,
		Logger:   logger,
		DB:       pool,
		Listener: listener,
		WebFS:    webFS,
	})

	logger.Info("starting", "addr", listener.Addr().String())
	return app.Run(ctx)
}

func newLogger(c config.Log) *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(c.Level))
	opts := &slog.HandlerOptions{Level: lvl}
	if c.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
