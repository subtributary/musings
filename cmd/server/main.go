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
		log.Fatalf("Error loading config: %v", err)
	}

	views, err := LoadViewFactory(cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", indexHandler(ctx, views, cfg))
	router.Get("/_static/*", staticHandler())
	router.Get("/*", contentHandler(views))

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
