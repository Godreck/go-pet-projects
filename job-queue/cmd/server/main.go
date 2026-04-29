package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/Godreck/go-pet-projects/job-queue/internal/http"
	"github.com/Godreck/go-pet-projects/job-queue/internal/job"
	"github.com/Godreck/go-pet-projects/job-queue/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st := store.New()
	manager := job.NewManager(st, 4, 128)
	manager.Start(ctx)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpapi.NewHandler(manager),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shotdownd error:", "error", err)
		}
	}()

	slog.Info("job-queue  server started", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server stoped with error", "error", err)
	}

	slog.Info("server stopped gracefully")
}
