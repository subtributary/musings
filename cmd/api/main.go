package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/web"
)

type Dependencies struct {
	ContentRoot *os.Root
}

func LoadDependencies() (*Dependencies, error) {
	deps := &Dependencies{}
	var err error

	deps.ContentRoot, err = os.OpenRoot(app.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
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
	router.Get("/files", filesGetHandler(deps))
	router.Get("/files/*", filesGetHandler(deps))

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

func filesGetHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := chi.URLParam(r, "*")
		response := web.NewResponse(w, r)

		files, err := fs.ReadDir(deps.ContentRoot.FS(), ".")
		if err != nil {
			response.ServerError(err)
			return
		}

		//
	}
}
