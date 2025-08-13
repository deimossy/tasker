package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/deimossy/tasker/internal/repository"
	"github.com/deimossy/tasker/internal/service"
	httptransport "github.com/deimossy/tasker/internal/transport/http"
)

func main() {
	logger := make(chan string, 100)
	defer close(logger)

	repo := repository.NewInMemoryTaskRepository()
	svc := service.NewTaskService(repo, logger)
	handler := httptransport.NewTaskHandler(svc)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited gracefully")
}
