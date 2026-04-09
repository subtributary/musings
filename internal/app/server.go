package app

import (
	"errors"
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

func NewServer(services Services, config Config) Server {
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

	return s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	locale := localization.LocaleFromContext(r.Context())

	root, err := os.OpenRoot(s.config.ContentPath)
	if err != nil {
		serveError(w, r, err)
	}
	defer func() { _ = root.Close() }()
	localizedFS := localization.NewLocalizedFS(root.FS(), locale)

	// If the specific file exists, serve it directly.
	if _, err = fs.Stat(localizedFS, path); err == nil {
		http.ServeFileFS(w, r, localizedFS, path)
		return
	}

	// todo: handle http head here please

	postData, err := s.services.PostParser.Parse(localizedFS, path+".md")
	if err != nil {
		serveError(w, r, err)
		return
	}

	if err = s.services.TemplateStore.Execute(w, "post", locale, postData); err != nil {
		serveError(w, r, err)
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	locale := localization.LocaleFromContext(r.Context())
	if err := s.services.TemplateStore.Execute(w, "index", locale, nil); err != nil {
		serveError(w, r, err)
	}
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	locale := localization.LocaleFromContext(r.Context())

	root, err := os.OpenRoot(s.config.GetStaticPath())
	if err != nil {
		serveError(w, r, err)
	}
	defer func() { _ = root.Close() }()
	localizedFS := localization.NewLocalizedFS(root.FS(), locale)

	http.ServeFileFS(w, r, localizedFS, path)
}

func serveError(w http.ResponseWriter, r *http.Request, err error) {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		log.Printf("file not found: %v", err)
		http.NotFound(w, r)
	} else {
		log.Printf("server error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
