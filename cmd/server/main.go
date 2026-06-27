package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	staticRoot, err := os.OpenRoot(StaticPath)
	if err != nil {
		log.Fatalf("Error: open static root: %v", err)
	}

	contentRoot, err := os.OpenRoot(ContentPath)
	if err != nil {
		log.Fatalf("Error: open content root: %v", err)
	}

	content, err := OpenContent(contentRoot, cfg.Locales)
	if err != nil {
		log.Fatalf("Error: load content: %v", err)
	}

	responder, err := NewResponder(cfg.LiveTemplates, cfg.Locales, contentRoot, staticRoot)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", indexHandler(responder, content))
	router.Get("/_static/*", fileHandler("/_static/", responder, staticRoot))
	router.Get("/*", contentHandler(responder, content))

	server := &http.Server{Addr: cfg.BindAddress, Handler: router}
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	log.Printf("Listening at %s\n", cfg.BindAddress)

	// Listen for CTRL+C or other interrupt.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server..")
	case err := <-serverErr:
		if err != nil {
			log.Printf("Error from server: %v", err)
		}
	}

	log.Println("Shutting down server...")
	stopCtx, stopStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopStop()
	err = server.Shutdown(stopCtx) // Shutdown server before closing deps.
	err = errors.Join(err, content.Close())
	err = errors.Join(err, staticRoot.Close())
	err = errors.Join(err, contentRoot.Close())
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Printf("Server stopped.")
}

func contentHandler(response Responder, content *Content) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPath, _ := strings.CutPrefix(r.URL.Path, "/")

		modTime, ok := content.ModTime(reqPath)

		// If ok, it is a regular file and not a post.
		if ok {
			response.File(w, r, content.root, reqPath)
			return
		}

		modTime, ok = content.ModTime(reqPath + ".md")
		if !ok {
			response.NotFound(w, r)
			return
		}

		post, err := content.GetPost(reqPath + ".md")
		if err != nil {
			response.NotFound(w, r)
		}

		response.View(w, r, "post", WithData(post, modTime))
	}
}

func fileHandler(prefix string, response Responder, root *os.Root) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.File(w, r, root, r.URL.Path)
	})
	return http.StripPrefix(prefix, handler).ServeHTTP
}

func indexHandler(response Responder, content *Content) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")
		results := content.Search(locale, query)
		response.View(w, r, "index", WithData(results, time.Now()))
	}
}
