package main

import (
	"context"
	"log"
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
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("job-queue server started on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped with error: %v", err)
	}

	log.Println("server stopped gracefully")
}
