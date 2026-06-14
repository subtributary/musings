package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

type Dependencies struct {
	ContentRoot *os.Root
	StaticRoot  *os.Root
	Posts       map[string]posts.Store // store := Posts[locale.Tag]
	Views       *ViewFactory
}

func LoadDependencies(cfg Config) (d Dependencies, err error) {
	d.ContentRoot, err = os.OpenRoot(ContentPath)
	if err != nil {
		err = fmt.Errorf("open content root: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.StaticRoot, err = os.OpenRoot(StaticPath)
	if err != nil {
		err = fmt.Errorf("open static root: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.Views, err = LoadViewFactory(cfg, d.StaticRoot)
	if err != nil {
		err = fmt.Errorf("load views: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.Posts = make(map[string]posts.Store)

	// Single "und" posts store if localization is disabled.
	if len(cfg.Locales) == 0 {
		store, err := posts.OpenStore(ContentPath)
		if err != nil {
			err = fmt.Errorf("open posts store: %w", err)
			return d, errors.Join(err, d.Close())
		}
		d.Posts["und"] = store
	}

	// Posts store per configured locale
	for _, locale := range cfg.Locales {
		locPath := filepath.Join(ContentPath, locale.Tag)
		store, err := posts.OpenStore(locPath)
		if err != nil {
			err = fmt.Errorf("open posts store for locale %s: %w", locale.Tag, err)
			return d, errors.Join(err, d.Close())
		}
		d.Posts[locale.Tag] = store
	}

	return d, nil
}

func (d *Dependencies) Close() (err error) {
	if d.ContentRoot != nil {
		err = errors.Join(err, d.ContentRoot.Close())
		d.ContentRoot = nil
	}

	if d.StaticRoot != nil {
		err = errors.Join(err, d.StaticRoot.Close())
		d.StaticRoot = nil
	}

	for loc, store := range d.Posts {
		err = errors.Join(err, store.Close())
		delete(d.Posts, loc)
	}

	return err
}

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
	router.Get("/", IndexHandler(deps.Views, deps.Posts))
	router.Get("/_static/*", StaticHandler(deps.Views, deps.StaticRoot))
	router.Get("/*", ContentHandler(deps.Views, deps.ContentRoot))

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
