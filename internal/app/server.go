package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	config   *Config
	router   *chi.Mux
	services *Services
}

func NewServer(services *Services, config *Config) (*Server, error) {
	s := &Server{
		config:   config,
		router:   chi.NewRouter(),
		services: services,
	}

	s.router.Use(middleware.Logger)
	s.router.Get("/", s.getIndex)
	s.router.Get("/_static/*", s.getStaticAsset)
	s.router.Get("/*", s.getPost)

	return s, nil
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) getIndex(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.

	store, err := s.services.TemplateProvider.Get()
	if err != nil {
		log.Printf("could not load templates: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = store.Execute(w, "index", locale, nil); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) getPost(w http.ResponseWriter, r *http.Request) {
	const locale = "" // I'll set this next update.

	// path is validated at the os boundary via `os.Root`.
	path := chi.URLParam(r, "*")

	filePath, err := s.services.MarkdownStore.Find(path, locale)
	if err != nil {
		log.Printf("error finding markdown file: %v", err)
		http.NotFound(w, r)
		return
	}

	fileData, err := s.services.MarkdownStore.Parse(filePath)
	if err != nil {
		log.Printf("error parsing markdown file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	store, err := s.services.TemplateProvider.Get()
	if err != nil {
		log.Printf("could not load templates: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = store.Execute(w, "post", locale, fileData); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) getStaticAsset(w http.ResponseWriter, r *http.Request) {
	staticDir := http.Dir(s.config.GetStaticPath())
	fs := http.StripPrefix("/_static/", http.FileServer(staticDir))
	fs.ServeHTTP(w, r)
}
