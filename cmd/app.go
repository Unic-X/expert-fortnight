package main

import (
	"context"
	"evently/internal/delivery/http/routes"
	"evently/internal/di"
	"log"
	"net/http"
	"time"
)

type Application struct {
	container *di.Container
	server    *http.Server
}

func NewApplication(container *di.Container) *Application {
	return &Application{
		container: container,
	}
}

func (a *Application) Start(ctx context.Context) {
	a.server = &http.Server{
		Addr:    ":8080",
		Handler: a.container.Server,
	}

	routes.AllRoutes(a.container.Server, a.container, a.container.JWTMiddleware)

	go func() {
		log.Println("Starting server on :8080")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()
}
