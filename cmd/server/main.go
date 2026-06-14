package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := LoadConfig()
	if err != nil {
		// LoadConfig already prints a friendly error message, so just return.
		return
	}

	views, err := LoadViewFactory(cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	handlers, err := NewHandlers(cfg, ctx, views)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer (func() { _ = handlers.Close() })()

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", handlers.IndexHandler())
	router.Get("/_static/*", handlers.StaticHandler())
	router.Get("/*", handlers.ContentHandler())

	log.Printf("Listening at %s\n", cfg.BindAddress)
	server := NewServer(cfg.BindAddress, router)
	server.Start()

	<-ctx.Done()
	stop()

	log.Println("Shutting down server...")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped.")
}
