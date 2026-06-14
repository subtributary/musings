package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		// LoadConfig already prints a friendly error message, so just return.
		return
	}

	deps, err := LoadDependencies(cfg)
	if err != nil {
		log.Fatalf("Error loading dependencies: %v", err)
	}
	defer (func() { _ = deps.Close() })()

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", IndexHandler(deps))
	router.Get("/_static/*", StaticHandler(deps))
	router.Get("/*", ContentHandler(deps))

	server := &http.Server{
		Addr:    cfg.BindAddress,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error from server: %v", err)
		}
	}()

	log.Printf("Listening at %s\n", cfg.BindAddress)

	// Listen for CTRL+C or other interrupt.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down server...")
	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownStop()
	err = server.Shutdown(shutdownCtx)
	err = errors.Join(err, deps.Close())
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Printf("Server stopped.")
}
