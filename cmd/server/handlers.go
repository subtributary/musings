package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
)

type Handlers struct {
	cfg         Config
	views       *ViewFactory
	ctx         context.Context
	contentRoot *os.Root
	staticRoot  *os.Root
}

func NewHandlers(cfg Config, ctx context.Context, views *ViewFactory) (*Handlers, error) {
	h := &Handlers{cfg: cfg, ctx: ctx, views: views}
	var err error

	h.contentRoot, err = os.OpenRoot(config.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %v", err)
	}

	h.staticRoot, err = os.OpenRoot(config.StaticPath)
	if err != nil {
		_ = h.contentRoot.Close()
		return nil, fmt.Errorf("open static root: %v", err)
	}

	return h, nil
}

func (h *Handlers) Close() error {
	contentErr := h.contentRoot.Close()
	staticErr := h.staticRoot.Close()
	return errors.Join(contentErr, staticErr)
}

func (h *Handlers) ContentHandler() http.HandlerFunc {
	fileHandler := h.FileHandler(h.contentRoot.FS())
	parser := posts.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimLeft(r.URL.Path, "/") + ".md"

		info, err := fs.Stat(h.contentRoot.FS(), filePath)
		switch {
		case os.IsNotExist(err):
			// It's not a Markdown file, so use the file server.
			fileHandler.ServeHTTP(w, r)
			return
		case errors.Is(err, fs.ErrInvalid):
			h.writeBadRequest(w)
			return
		case err != nil:
			h.writeServerError(w, err)
			return
		}

		content, err := parser.ParseFile(h.contentRoot.FS(), filePath)
		if err != nil {
			h.writeServerError(w, err)
			return
		}

		err = h.views.CreateAndServe(w, "post",
			WithData(content),
			WithDataModified(info.ModTime()),
			WithLocale(localization.LocaleFromContext(r.Context())),
			WithPath(chi.RouteContext(r.Context()).RoutePath),
		)
		if err != nil {
			h.writeServerError(w, fmt.Errorf("serve view: %v", err))
		}
	}
}

func (h *Handlers) FileHandler(root fs.FS) http.HandlerFunc {
	wrapped := http.FileServerFS(root)

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Stat(root, r.URL.Path)
		if errors.Is(err, fs.ErrNotExist) {
			h.writeNotFound(w, r)
		}

		wrapped.ServeHTTP(w, r)
	}
}

func (h *Handlers) IndexHandler() http.HandlerFunc {
	locales := h.cfg.Locales

	// Even with no locales we still need a store, so use the Und locale.
	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}

	stores := make(map[string]*posts.Store)
	for _, locale := range locales {
		locPath := config.LocalizedContentPath(locale)
		store, err := posts.OpenStore(h.ctx, locPath)
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
	handler = h.FileHandler(h.staticRoot.FS())
	handler = http.StripPrefix("/_static/", handler)
	return handler.ServeHTTP
}

func (h *Handlers) writeBadRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (h *Handlers) writeNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	err := h.views.CreateAndServe(w, "404",
		WithLocale(localization.LocaleFromContext(r.Context())),
		WithPath(chi.RouteContext(r.Context()).RoutePath),
	)
	if err != nil {
		h.writeServerError(w, err)
	}
}

func (h *Handlers) writeServerError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
