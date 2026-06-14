package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

func main() {
	// Context so we can gracefully shut down from CTRL+C or other interrupt.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	roots, err := OpenRoots()
	if err != nil {
		log.Fatalf("Error opening roots: %v", err)
	}
	defer (func() { _ = roots.Close() })()

	cfg, err := LoadConfig()
	if err != nil {
		// LoadConfig already prints a friendly error message, so just return.
		return
	}

	deps, err := loadDependencies(ctx, roots, cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", IndexHandler(deps.Views, deps.PostStores))
	router.Get("/_static/*", StaticHandler(deps.Views, roots.Static))
	router.Get("/*", ContentHandler(deps.Views, roots.Content))

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

type Deps struct {
	PostStores map[string]*posts.Store
	Views      *ViewFactory
}

func loadDependencies(ctx context.Context, roots Roots, cfg Config) (Deps, error) {
	postStores := make(map[string]*posts.Store)

	// Single "und" post store if localization is disabled.
	if len(cfg.Locales) == 0 {
		store, err := posts.OpenStore(ctx, ContentPath)
		if err != nil {
			return Deps{}, fmt.Errorf("open posts store: %w", err)
		}
		postStores["und"] = store
	}

	// Post store per configured locale
	for _, locale := range cfg.Locales {
		locPath := filepath.Join(ContentPath, locale.Tag)
		store, err := posts.OpenStore(ctx, locPath)
		if err != nil {
			return Deps{}, fmt.Errorf("open posts store for locale %s: %w", locale.Tag, err)
		}
		postStores[locale.Tag] = store
	}

	views, err := LoadViewFactory(cfg, roots.Static)
	if err != nil {
		return Deps{}, fmt.Errorf("load view factory: %w", err)
	}

	return Deps{
		PostStores: postStores,
		Views:      views,
	}, nil
}
