package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
)

type Server struct {
	config   Config
	router   *chi.Mux
	services Services
}

func NewServer(services Services, config Config) *Server {
	s := Server{
		config:   config,
		router:   chi.NewRouter(),
		services: services,
	}

	s.router.Use(middleware.Logger)
	s.router.Use(localization.LocalizedRoute(s.config.Locales))
	s.router.Get("/", s.handleIndex)
	s.router.Get("/_static/*", s.handleStatic)
	s.router.Get("/*", s.handleContent)

	return &s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	if err := serveFile(w, r, s.config.ContentPath, path); err != nil {
		if err = s.servePost(w, r, path); err != nil {
			writeError(w, err)
		}
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.writeTemplate(w, r, "index", nil); err != nil {
		writeError(w, err)
	}
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	if err := serveFile(w, r, s.config.GetStaticPath(), path); err != nil {
		writeError(w, err)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, rootPath string, name string) error {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return err
	}
	defer func() { _ = root.Close() }()
	locale := localization.LocaleFromContext(r.Context())
	localizedFS := localization.NewLocalizedFS(root.FS(), locale)

	// If the file exists, then serve it.
	if _, err = fs.Stat(localizedFS, name); err == nil {
		http.ServeFileFS(w, r, localizedFS, name)
	}
	return err
}

func (s *Server) servePost(w http.ResponseWriter, r *http.Request, path string) error {
	root, err := os.OpenRoot(s.config.ContentPath)
	if err != nil {
		return err
	}
	defer func() { _ = root.Close() }()
	locale := localization.LocaleFromContext(r.Context())
	localizedFS := localization.NewLocalizedFS(root.FS(), locale)

	// todo: handle http head here please

	path += ".md"
	data, err := s.services.PostParser.Parse(localizedFS, path)
	if err != nil {
		return err
	}
	return s.writeTemplate(w, r, "post", data)
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

	tmpl, err := s.services.TemplateStore.Lookup(name, locale)
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
