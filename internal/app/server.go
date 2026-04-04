package app

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/store"
)

type Server struct {
	config   Config
	router   *chi.Mux
	services Services
}

func NewServer(services Services, config Config) (Server, error) {
	s := Server{
		config:   config,
		router:   chi.NewRouter(),
		services: services,
	}

	s.router.Use(middleware.Logger)
	s.router.Get("/", s.handleIndexGet)
	s.router.Get("/_static/*", s.handleStaticGet)
	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if !s.handleContentGet(w, r) {
			s.handlePostGet(w, r)
		}
	})

	return s, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

// handleContentGet handles GET requests for content.
// It returns true if it handles the request.
func (s *Server) handleContentGet(w http.ResponseWriter, r *http.Request) bool {
	const locale = "" // I'll set this next update.
	path := chi.URLParam(r, "*")

	foundPath, err := s.services.ContentStore.Find(path, locale)
	if err != nil || foundPath == "" {
		return false
	}

	http.ServeFile(w, r, filepath.Join(s.config.ContentPath, foundPath))
	return true
}

func (s *Server) handleIndexGet(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.
	s.writeTemplate(w, "index", locale, nil)
}

func (s *Server) handlePostGet(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.
	path := chi.URLParam(r, "*") + ".md"

	fileData, err := s.services.PostsStore.Parse(path, locale)
	if errors.Is(err, &store.NotFoundError{}) {
		log.Printf("error finding markdown file: %v", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("error parsing markdown file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.writeTemplate(w, "post", locale, fileData)
}

func (s *Server) handleStaticGet(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.
	path := chi.URLParam(r, "*")

	foundPath, err := s.services.StaticStore.Find(path, locale)
	if errors.Is(err, &store.NotFoundError{}) {
		log.Printf("error finding static file: %v", err)
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("error parsing markdown file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filepath.Join(s.config.GetStaticPath(), foundPath))
}

func (s *Server) writeTemplate(w http.ResponseWriter, name string, locale string, data any) {
	templatesStore, err := s.services.TemplateProvider.Get()
	if err != nil {
		log.Printf("could not load templates: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl := templatesStore.Lookup(name, locale)
	if tmpl == nil {
		// This isn't a 404 because all templates referenced should be present.
		log.Printf("could not locate template: %q", name)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
