package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/files"
)

type Server struct {
	config    *Config
	router    *chi.Mux
	templates *files.TemplateStore
	markdown  *files.MarkdownStore
}

func NewServer(config *Config) (*Server, error) {
	templates := files.NewTemplateStore()
	if err := templates.Load(config.GetTemplatesPath()); err != nil {
		return nil, fmt.Errorf("could not load templates: %w", err)
	}

	markdown, err := files.NewMarkdownStore(config.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("could not load markdown store: %w", err)
	}

	s := &Server{
		config:    config,
		router:    chi.NewRouter(),
		templates: templates,
		markdown:  markdown,
	}

	s.router.Use(middleware.Logger)
	s.router.Get("/", s.getIndex)
	s.router.Get("/*", s.getPost)

	return s, nil
}

func (s *Server) Dispose() {
	s.markdown.Dispose()
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) getPost(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.

	// path is validated at the os boundary via `os.Root`.
	path := chi.URLParam(r, "*")

	filePath, err := s.markdown.Find(path, locale)
	if err != nil {
		log.Printf("error finding markdown file: %v", err)
		http.NotFound(w, r)
		return
	}

	fileData, err := s.markdown.Parse(filePath)
	if err != nil {
		log.Printf("error parsing markdown file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	template := s.templates.Lookup("post", locale)
	if template == nil {
		log.Printf("no post template found for locale %s", locale)
		http.NotFound(w, r)
		return
	}

	err = template.Execute(w, fileData)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
