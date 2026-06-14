package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

type Handlers struct {
	posts       map[string]*posts.Store // store := posts[locale.Tag]
	views       *ViewFactory
	contentRoot *os.Root
	staticRoot  *os.Root
}

func NewHandlers(cfg Config, ctx context.Context, views *ViewFactory) (*Handlers, error) {
	h := &Handlers{views: views}
	var err error

	h.contentRoot, err = os.OpenRoot(ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
	}

	h.staticRoot, err = os.OpenRoot(StaticPath)
	if err != nil {
		_ = h.contentRoot.Close()
		return nil, fmt.Errorf("open static root: %w", err)
	}

	postsStore, err := newPostsStore(ctx, cfg.Locales)
	if err != nil {
		_ = h.contentRoot.Close()
		_ = h.staticRoot.Close()
		return nil, fmt.Errorf("open posts store: %w", err)
	}
	h.posts = postsStore

	return h, nil
}

func newPostsStore(ctx context.Context, locales []localization.Locale) (map[string]*posts.Store, error) {
	// Even with no locales we still need a store, so use the Und locale.
	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}

	stores := make(map[string]*posts.Store)

	for _, locale := range locales {
		locPath := ContentPath
		if locale != localization.UndLocale {
			locPath = filepath.Join(locPath, locale.Tag)
		}

		store, err := posts.OpenStore(ctx, locPath)
		if err != nil {
			return nil, fmt.Errorf("could not load store: %w", err)
		}

		stores[locale.Tag] = store
	}

	return stores, nil
}

func (h *Handlers) Close() error {
	contentErr := h.contentRoot.Close()
	staticErr := h.staticRoot.Close()
	return errors.Join(contentErr, staticErr)
}

func (h *Handlers) ContentHandler() http.HandlerFunc {
	parser := posts.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")

		// Try opening as regular file, then try opening as post.
		isPost := false
		file, err := h.contentRoot.Open(name)
		if err != nil {
			file, err = h.contentRoot.Open(name + ".md")
			if err != nil {
				h.writeNotFound(w, r, name)
				return
			}
			isPost = true
		}
		defer func() { _ = file.Close() }()

		info, err := file.Stat()
		if err != nil {
			h.writeServerError(w, fmt.Errorf("stat file: %w", err))
			return
		}

		if info.IsDir() {
			h.writeNotFound(w, r, name)
			return
		}

		if !isPost {
			http.ServeContent(w, r, info.Name(), info.ModTime(), file)
			return
		}

		content, err := io.ReadAll(file)
		if err != nil {
			h.writeServerError(w, fmt.Errorf("read file: %w", err))
			return
		}

		post, err := parser.ParseContent(content)
		if err != nil {
			h.writeServerError(w, fmt.Errorf("parse post: %w", err))
			return
		}

		err = h.views.CreateAndServe(w, "post",
			WithData(post),
			WithDataModified(info.ModTime()),
			WithLocale(localization.LocaleFromContext(r.Context())),
			WithPath(chi.RouteContext(r.Context()).RoutePath),
		)
		if err != nil {
			h.writeServerError(w, fmt.Errorf("serve view: %w", err))
		}
	}
}

func (h *Handlers) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := localization.LocaleFromContext(r.Context())
		query := r.URL.Query().Get("q")

		store, ok := h.posts[locale.Tag]
		if !ok {
			// The set of configured locales should never change,
			// so we should never get this error.
			h.writeServerError(w, errors.New("posts store missing for locale"))
		}

		var results []posts.IndexedPost
		for result := range store.Search(query) {
			if locale != localization.UndLocale {
				result.Path = "/" + path.Join(locale.Tag, result.Path)
			}
			results = append(results, result)
		}

		err := h.views.CreateAndServe(w, "index",
			WithData(results),
			WithDataModified(time.Now()),
			WithLocale(localization.LocaleFromContext(r.Context())),
			WithPath(chi.RouteContext(r.Context()).RoutePath),
		)
		if err != nil {
			h.writeServerError(w, fmt.Errorf("serve view: %v", err))
		}
	}
}

func (h *Handlers) StaticHandler() http.HandlerFunc {
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path

		file, err := h.staticRoot.Open(name)
		if err != nil {
			h.writeNotFound(w, r, name)
			return
		}
		defer (func() { _ = file.Close() })()

		info, err := file.Stat()
		if err != nil {
			h.writeServerError(w, fmt.Errorf("stat file: %w", err))
			return
		}

		if info.IsDir() {
			h.writeNotFound(w, r, name)
			return
		}

		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	})

	handler = http.StripPrefix("/_static/", handler)
	return handler.ServeHTTP
}

func (h *Handlers) writeNotFound(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotFound)

	err := h.views.CreateAndServe(w, "404",
		WithLocale(localization.LocaleFromContext(r.Context())),
		WithPath(name),
	)
	if err != nil {
		log.Printf("Error serving not found response: %v", err)
	}
}

func (h *Handlers) writeServerError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
