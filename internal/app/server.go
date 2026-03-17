package app

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/content"
)

type Server struct {
	config    *Config
	router    *chi.Mux
	templates *content.TemplateStore
	markdown  *content.MarkdownStore
}

func NewServer(config *Config) (*Server, error) {
	s := &Server{
		config: config,
	}

	s.templates = content.NewTemplateStore()
	if err := s.templates.Load(config.GetTemplatesPath()); err != nil {
		return nil, fmt.Errorf("could not load templates: %w", err)
	}

	s.markdown = content.NewMarkdownStore(config.ContentPath)

	s.router = chi.NewRouter()
	s.router.Use(middleware.Logger)
	s.router.Get("/", s.getIndex)
	s.router.Get("/*", s.getPost)

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

func (s *Server) getPost(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.

	path := chi.URLParam(r, "*")
	if !fs.ValidPath(path) {
		http.NotFound(w, r)
		return
	}

	filePath, err := s.markdown.Find(path, locale)
	if err != nil {
		log.Printf("error finding markdown file: %v", err)
		http.NotFound(w, r)
		return
	}

	if err := s.markdown.Render(w, filePath); err != nil {
		log.Printf("error rendering markdown file: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
