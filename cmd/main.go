package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"evently/internal/di"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize dependency injection container
	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize container: %v", err)
	}
	defer container.Pool.Close()

	// Create application
	app := NewApplication(container)
	app.SetupRoutes()

	// Start server
	app.Start(ctx)

	<-ctx.Done()
	log.Println("shutdown signal received, closing resources...")
}
