package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/web"
)

//go:embed all:web
var webFiles embed.FS

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		// LoadConfig already prints a friendly error message, so just return.
		return
	}

	templates, err := LoadTemplates()
	if err != nil {
		log.Fatalf("Error: load templates: %v", err)
	}

	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", indexHandler(templates, cfg.Locales))
	router.Handle("/_css/*", http.FileServerFS(webFS))
	router.Handle("/_fonts/*", http.FileServerFS(webFS))
	router.Handle("/_images/*", http.FileServerFS(webFS))
	router.Get("/*", contentHandler(templates, cfg.Locales))

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
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Printf("Server stopped.")
}

func contentHandler(templates Templates, locales []localization.Locale) http.HandlerFunc {
	contentRoot, err := os.OpenRoot(ContentPath)
	if err != nil {
		log.Fatalf("Error: open content root: %v", err)
	}

	modTime := func(name string) (time.Time, bool) {
		info, err := contentRoot.Stat(name)
		if err != nil {
			return time.Time{}, false
		}
		return info.ModTime(), true
	}
	parser := posts.NewParser(modTime)

	vmFactory, err := NewModelViewFactory(locales)
	if err != nil {
		log.Fatalf("Error: create view model factory: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := web.NewResponse(w, r)
		filePath, _ := strings.CutPrefix(r.URL.Path, "/")

		// If no error, it is a regular file and not a post.
		if _, err := contentRoot.Stat(filePath); err == nil {
			response.File(contentRoot, filePath)
			return
		}

		if _, err := contentRoot.Stat(filePath + ".md"); err != nil {
			response.NotFound(templates.Err404, vmFactory.Create(r, nil))
			return
		}

		post, err := parser.ParseFile(contentRoot.FS(), filePath+".md")
		if err != nil {
			response.ServerError(err)
			return
		}

		response.View(templates.Post, vmFactory.Create(r, post))
	}
}

type SearchResults struct {
	Query   string
	Results []posts.IndexedPost
}

func indexHandler(templates Templates, locales []localization.Locale) http.HandlerFunc {
	vmFactory, err := NewModelViewFactory(locales)
	if err != nil {
		log.Fatalf("Error: create view model factory: %v", err)
	}

	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}

	indexes := make(map[string]posts.AutoIndex, len(locales))
	for _, loc := range locales {
		indexRoot := path.Join(ContentPath, loc.Tag)
		index, err := posts.NewAutoIndex(indexRoot)
		if err != nil {
			log.Fatalf("Error: load indexes: %v", err)
		}
		indexes[loc.Tag] = index
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqLocale := localization.LocaleFromContext(r.Context())

		index := indexes[reqLocale.Tag]
		query := r.URL.Query().Get("q")
		results := slices.Collect(index.Search(query))

		response := web.NewResponse(w, r)
		response.View(templates.Index, vmFactory.Create(r, SearchResults{
			Query:   query,
			Results: results,
		}))
	}
}
