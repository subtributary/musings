package main

import (
	"errors"
	"fmt"
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

func ContentHandler(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		locPath := chi.RouteContext(r.Context()).RoutePath
		locPath, _ = strings.CutPrefix(locPath, "/")

		store := deps.Posts[locale.Tag]
		post, err := store.Post(locPath + ".md")

		if err != nil {
			serveFile(w, r, deps.Views, deps.ContentRoot, r.URL.Path)
			return
		}

		err = deps.Views.CreateAndServe(w, "post",
			WithData(post),
			WithDataModified(post.Modified),
			WithLocale(locale),
			WithPath(locPath),
		)
		if err != nil {
			writeServerError(w, fmt.Errorf("serve view: %w", err))
		}
	}
}

func IndexHandler(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")

		locStore, ok := deps.Posts[locale.Tag]
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

		err := deps.Views.CreateAndServe(w, "index",
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

func StaticHandler(deps Dependencies) http.HandlerFunc {
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, r, deps.Views, deps.StaticRoot, r.URL.Path)
	})

	handler = http.StripPrefix("/_static/", handler)
	return handler.ServeHTTP
}

// serveFile safely serves the file with the name if it exists.
func serveFile(
	w http.ResponseWriter, r *http.Request, views *ViewFactory,
	root *os.Root, name string,
) {
	file, err := root.Open(name)
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
}

func writeNotFound(
	w http.ResponseWriter, r *http.Request, views *ViewFactory,
	name string,
) {
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
