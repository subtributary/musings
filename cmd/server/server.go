package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
)

type Server struct {
	config   Config
	router   *chi.Mux
	services Services
}

func NewServer(services Services, cfg Config) *Server {
	s := Server{
		config:   cfg,
		router:   chi.NewRouter(),
		services: services,
	}

	s.router.Use(middleware.Logger)
	s.router.Use(localization.LocalizedRoute(s.config.Locales))
	s.router.Get("/", s.handleIndex)
	s.router.Handle("/_static/*", http.StripPrefix("/_static/",
		http.FileServer(http.Dir(config.StaticPath)),
	))
	s.router.Get("/*", s.handleContent)

	return &s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*") + ".md"

	root, err := os.OpenRoot(config.ContentPath)
	if err != nil {
		writeError(w, err)
		return
	}
	defer func() { _ = root.Close() }()

	info, err := fs.Stat(root.FS(), path)
	if err != nil {
		// It's not a markdown file, so use the Go file server.
		fileServer := http.FileServer(http.Dir(config.ContentPath))
		fileServer.ServeHTTP(w, r)
		return
	}
	modTime := info.ModTime().UTC().Truncate(time.Second)

	if t, err := http.ParseTime(r.Header.Get("If-Modified-Since")); err == nil {
		log.Printf("modTime: %v, t: %v", modTime, t)
		if !modTime.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Last-Modified", modTime.Format(http.TimeFormat))

	if r.Method != "HEAD" {
		data, err := s.services.PostParser.ParseFile(root.FS(), path)
		if err != nil {
			writeError(w, err)
		}
		if err = s.writeTemplate(w, r, "post", data); err != nil {
			writeError(w, err)
		}
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	locale := localization.LocaleFromContext(r.Context())
	index := s.services.PostIndexes[locale]
	results := slices.Collect(index.Search(query))
	if err := s.writeTemplate(w, r, "index", results); err != nil {
		writeError(w, err)
	}
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

func (s *Server) writeTemplate(w http.ResponseWriter, r *http.Request, name string, data any) error {
	locale := localization.LocaleFromContext(r.Context())
	path := chi.RouteContext(r.Context()).RoutePath
	viewModel := NewViewModel(ViewModelParams{
		CurrentLocale:    locale,
		SupportedLocales: s.config.Locales,
		CurrentPath:      path,
		Data:             data,
	})

	tmpl, err := s.services.TemplateStore.Lookup(name)
	if err != nil {
		return fmt.Errorf("template %q not found: %w", name, err)
	}

	// Write to a buffer so that errors do not leave the template partially written.
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, viewModel); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}
