package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

func currentRoutePath(r *http.Request) string {
	return chi.URLParam(r, "*")
}

func contentHandler(views *ViewFactory) http.HandlerFunc {
	root, err := os.OpenRoot(config.ContentPath)
	if err != nil {
		log.Fatalf("error opening content root: %v", err)
	}

	fallback := http.FileServerFS(root.FS())
	parser := posts.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimLeft(r.URL.Path, "/") + ".md"

		info, err := fs.Stat(root.FS(), filePath)
		if errors.Is(err, fs.ErrNotExist) {
			// It's not a Markdown file, so use the Go file server.
			fallback.ServeHTTP(w, r)
			return
		} else if err != nil {
			writeError(w, err)
			return
		}

		content, err := parser.ParseFile(root.FS(), filePath)
		if err != nil {
			writeError(w, err)
			return
		}

		err = views.Serve(w, r, "post",
			WithData(content),
			WithDataModified(info.ModTime()),
		)
		if err != nil {
			writeError(w, err)
			return
		}
	}
}

func indexHandler(ctx context.Context, views *ViewFactory, cfg Config) http.HandlerFunc {
	locales := cfg.Locales

	// Even with no locales we still need a store, so use the Und locale.
	if len(cfg.Locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}

	stores := make(map[string]*posts.Store)
	for _, locale := range locales {
		locPath := config.LocalizedContentPath(locale)
		store, err := posts.OpenStore(ctx, locPath)
		if err != nil {
			log.Fatalf("could not load store: %v", err)
		}
		stores[locale.Tag] = store
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")

		var results []posts.IndexedPost
		for result := range stores[locale.Tag].Search(query) {
			if locale != localization.UndLocale {
				tag := strings.ToLower(locale.Tag)
				result.Path = "/" + path.Join(tag, result.Path)
			}
			results = append(results, result)
		}

		err := views.Serve(w, r, "index",
			WithData(results),
		)
		if err != nil {
			writeError(w, err)
			return
		}
	}
}

func staticHandler() http.HandlerFunc {
	var handler http.Handler
	handler = http.FileServer(http.Dir(config.StaticPath))
	handler = http.StripPrefix("/_static/", handler)
	return handler.ServeHTTP
}

func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, fs.ErrNotExist) {
		log.Printf("file not found: %v", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		log.Printf("server error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
