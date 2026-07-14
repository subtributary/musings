package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/localization"
)

//go:embed all:templates
var templateFiles embed.FS

//go:embed all:web
var webFiles embed.FS

type Dependencies struct {
	ContentRoot *os.Root
	WebFS       fs.FS

	IndexTemplate *template.Template
}

func LoadDependencies() (*Dependencies, error) {
	deps := &Dependencies{}
	var err error

	deps.ContentRoot, err = os.OpenRoot(app.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
	}

	deps.WebFS, err = fs.Sub(webFiles, "web")
	if err != nil {
		return nil, fmt.Errorf("sub embedded web files: %w", err)
	}

	templateFS, err := fs.Sub(templateFiles, "templates")
	if err != nil {
		return nil, fmt.Errorf("sub embedded template dir: %w", err)
	}

	deps.IndexTemplate, err = template.ParseFS(templateFS, "index.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load index template: %w", err)
	}

	return deps, nil
}

func (d *Dependencies) Close() error {
	err := d.ContentRoot.Close()
	return err
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			PrintUsage()
		} else if argsErr, ok := app.AsArgsError(err); ok {
			PrintArgsErr(argsErr)
		} else {
			log.Fatalf("Error: load config: %v", err)
		}
		return
	}

	deps, err := LoadDependencies()
	if err != nil {
		log.Fatalf("Error: load deps: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Handle("/*", http.FileServerFS(deps.WebFS))

	server := app.StartServer(cfg.BindAddress, router)
	log.Printf("Listening at %s\n", cfg.BindAddress)
	if err = server.Wait(); err != nil {
		log.Printf("Error from server: %v", err)
	}

	log.Println("Shutting down server.")
	err = server.Close() // Shutdown server before closing deps.
	err = errors.Join(err, deps.Close())
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Printf("Server stopped.")
}
