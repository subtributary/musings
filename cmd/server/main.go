package main

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"golang.org/x/text/language"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	localization.InitTranslations()
	views, err := LoadViewFactory(cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", indexHandler(views, cfg))
	router.Get("/_static/*", staticHandler())
	router.Get("/*", contentHandler(views))

	log.Printf("Listening at %s\n", cfg.BindAddress)
	log.Fatal(http.ListenAndServe(cfg.BindAddress, router))
}

func currentRoutePath(r *http.Request) string {
	return chi.URLParam(r, "*")
}

func contentHandler(views *ViewFactory) http.HandlerFunc {
	fallback := http.FileServer(http.Dir(config.ContentPath))
	parser := posts.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		root, err := os.OpenRoot(config.ContentPath)
		if err != nil {
			writeError(w, err)
			return
		}
		defer func() { _ = root.Close() }()

		path := currentRoutePath(r) + ".md"

		info, err := fs.Stat(root.FS(), path)
		if err != nil {
			// It's not a Markdown file, so use the Go file server.
			fallback.ServeHTTP(w, r)
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

func indexHandler(views *ViewFactory, cfg Config) http.HandlerFunc {
	locales := cfg.Locales
	if len(locales) == 0 {
		locales = []language.Tag{language.Und}
	}

	indexes := make(map[language.Tag]posts.Index)
	for _, tag := range locales {
		locale := tag.String()
		index, err := posts.LoadIndex(config.DataPath, locale)
		if err != nil {
			log.Fatalf("could not load index: %v", err)
		}
		indexes[tag] = index
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")
		results := slices.Collect(indexes[locale].Search(query))

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
