package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/templates"
)

type Server struct {
	config    *Config
	router    *chi.Mux
	templates *templates.Store
}

func NewServer(config *Config) (*Server, error) {
	templatesStore, err := templates.New(config.GetTemplatesPath())
	if err != nil {
		return nil, fmt.Errorf("error creating templates store: %w", err)
	}

	s := &Server{
		config:    config,
		router:    chi.NewRouter(),
		templates: templatesStore,
	}

	s.router.Use(middleware.Logger)

	s.router.Get("/", s.getIndex)

	return s, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) getIndex(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.
	template := s.templates.Lookup("index", locale)
	if template == nil {
		log.Printf("no index template found for locale %s", locale)
		http.NotFound(w, r)
		return
	}

	err := template.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
}

//
