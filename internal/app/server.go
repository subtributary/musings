package app

import (
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func NewServer(config *Config) *http.Server {
	imagesDir := filepath.Join(config.ContentDir, "images")

	mux := http.NewServeMux()
	mux.Handle("/_content/images/", http.StripPrefix("/_content/images/", http.FileServer(http.Dir(imagesDir))))

	// Routes (todo)
	// GET /_api/ - api
	// GET /_assets/ - dist assets directory
	// GET /_content/images - content images
	// GET /loc/posts/1/slug - content markdown
	// GET /loc/page - posts as custom routes

	handler := loggingMiddleware(mux)

	return &http.Server{
		Addr:         config.WebEndpoint,
		Handler:      handler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
