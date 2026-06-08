package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"golang.org/x/text/language"
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
		path := strings.TrimLeft(r.URL.Path, "/") + ".md"

		info, err := fs.Stat(root.FS(), path)
		if errors.Is(err, fs.ErrNotExist) {
			// It's not a Markdown file, so use the Go file server.
			fallback.ServeHTTP(w, r)
			return
		} else if err != nil {
			writeError(w, err)
			return
		}

		content, err := parser.ParseFile(root.FS(), path)
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
	if len(locales) == 0 {
		locales = []language.Tag{language.Und}
	}

	stores := make(map[language.Tag]*posts.Store)
	for _, tag := range locales {
		path := config.LocalizedContentPath(tag)
		store, err := posts.OpenStore(ctx, path)
		if err != nil {
			log.Fatalf("could not load store: %v", err)
		}
		stores[tag] = store
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")
		results := slices.Collect(stores[locale].Search(query))

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
