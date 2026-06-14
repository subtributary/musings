package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

func ContentHandler(views *ViewFactory, contentRoot *os.Root) http.HandlerFunc {
	parser := posts.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")

		// Try opening as regular file, then try opening as post.
		isPost := false
		file, err := contentRoot.Open(name)
		if err != nil {
			file, err = contentRoot.Open(name + ".md")
			if err != nil {
				writeNotFound(w, r, views, name)
				return
			}
			isPost = true
		}
		defer func() { _ = file.Close() }()

		info, err := file.Stat()
		if err != nil {
			writeServerError(w, fmt.Errorf("stat file: %w", err))
			return
		}

		if info.IsDir() {
			writeNotFound(w, r, views, name)
			return
		}

		if !isPost {
			http.ServeContent(w, r, info.Name(), info.ModTime(), file)
			return
		}

		content, err := io.ReadAll(file)
		if err != nil {
			writeServerError(w, fmt.Errorf("read file: %w", err))
			return
		}

		post, err := parser.ParseContent(content)
		if err != nil {
			writeServerError(w, fmt.Errorf("parse post: %w", err))
			return
		}

		err = views.CreateAndServe(w, "post",
			WithData(post),
			WithDataModified(info.ModTime()),
			WithLocale(localization.LocaleFromContext(r.Context())),
			WithPath(chi.RouteContext(r.Context()).RoutePath),
		)
		if err != nil {
			writeServerError(w, fmt.Errorf("serve view: %w", err))
		}
	}
}

func IndexHandler(views *ViewFactory, stores map[string]posts.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")

		locStore, ok := stores[locale.Tag]
		if !ok {
			// The set of configured locales should never change,
			// so we should never get this error.
			writeServerError(w, errors.New("posts store missing for locale"))
		}

		var results []posts.IndexedPost
		for result := range locStore.Search(query) {
			if locale != localization.UndLocale {
				result.Path = "/" + path.Join(locale.Tag, result.Path)
			}
			results = append(results, result)
		}

		err := views.CreateAndServe(w, "index",
			WithData(results),
			WithDataModified(time.Now()),
			WithLocale(localization.LocaleFromContext(r.Context())),
			WithPath(chi.RouteContext(r.Context()).RoutePath),
		)
		if err != nil {
			writeServerError(w, fmt.Errorf("serve view: %v", err))
		}
	}
}

func StaticHandler(views *ViewFactory, staticRoot *os.Root) http.HandlerFunc {
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path

		file, err := staticRoot.Open(name)
		if err != nil {
			writeNotFound(w, r, views, name)
			return
		}
		defer (func() { _ = file.Close() })()

		info, err := file.Stat()
		if err != nil {
			writeServerError(w, fmt.Errorf("stat file: %w", err))
			return
		}

		if info.IsDir() {
			writeNotFound(w, r, views, name)
			return
		}

		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	})

	handler = http.StripPrefix("/_static/", handler)
	return handler.ServeHTTP
}

func writeNotFound(w http.ResponseWriter, r *http.Request, views *ViewFactory, name string) {
	w.WriteHeader(http.StatusNotFound)

	err := views.CreateAndServe(w, "404",
		WithLocale(localization.LocaleFromContext(r.Context())),
		WithPath(name),
	)
	if err != nil {
		log.Printf("Error serving not found response: %v", err)
	}
}

func writeServerError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
