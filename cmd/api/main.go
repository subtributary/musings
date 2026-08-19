package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/web"
)

type Dependencies struct {
	ContentRoot *os.Root
	DataRoot    *os.Root
}

func LoadDependencies() (*Dependencies, error) {
	deps := &Dependencies{}
	var err error

	deps.ContentRoot, err = os.OpenRoot(app.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
	}

	deps.DataRoot, err = os.OpenRoot(app.DataPath)
	if err != nil {
		return nil, fmt.Errorf("open data root: %w", err)
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
	router.Use(cors.AllowAll().Handler)
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Get("/content", contentGetHandler(deps))
	router.Get("/content/*", contentGetHandler(deps))
	router.Get("/data/*", dataGetHandler(deps))
	router.Get("/index", indexGetHandler(deps))

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

func contentGetHandler(deps *Dependencies) http.HandlerFunc {
	type ResponseItem struct {
		IsDir bool   `json:"is_dir"`
		Name  string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "*")
		if name == "" {
			name = "."
		}

		response := web.NewJSONResponse(w, r)

		info, err := deps.ContentRoot.Stat(name)
		if err != nil {
			response.NotFound()
		}

		if !info.IsDir() {
			response.File(deps.ContentRoot, name)
			return
		}

		files, err := fs.ReadDir(deps.ContentRoot.FS(), name)
		if err != nil {
			response.ServerError(err)
			return
		}

		items := make([]ResponseItem, len(files))
		for i, f := range files {
			items[i] = ResponseItem{
				IsDir: f.IsDir(),
				Name:  f.Name(),
			}
		}

		response.Okay(items)
	}
}

func dataGetHandler(deps *Dependencies) http.HandlerFunc {
	whitelist := []string{
		"config.json",
		"translations.json",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "*")
		if name == "" {
			name = "."
		}

		response := web.NewJSONResponse(w, r)

		if !slices.Contains(whitelist, name) {
			response.NotFound()
		}

		response.File(deps.DataRoot, name)
	}
}

func indexGetHandler(deps *Dependencies) http.HandlerFunc {
	index, err := posts.NewAutoIndex(app.ContentPath)
	if err != nil {
		log.Fatalf("Error: load index: %v", err)
	}

	type ResponseModel struct {
		Query   string               `json:"query"`
		Results []*posts.IndexedPost `json:"results"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		model := ResponseModel{
			Query:   r.URL.Query().Get("q"),
			Results: make([]*posts.IndexedPost, 0),
		}

		for result := range index.Search(model.Query) {
			if result.Thumbnail != "" {
				postPath := path.Dir(result.Path)
				thumbnail := web.VersionURL(deps.ContentRoot, postPath, result.Thumbnail)
				result.Thumbnail = path.Join(postPath, thumbnail)
			}

			model.Results = append(model.Results, result)
		}

		response := web.NewJSONResponse(w, r)
		response.Okay(model)
	}
}
