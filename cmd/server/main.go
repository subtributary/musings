package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/web"
)

//go:embed all:templates
var templateFiles embed.FS

//go:embed all:web
var webFiles embed.FS

// Dependencies are shared by multiple routes or require closing.
type Dependencies struct {
	ContentRoot *os.Root
	ViewModels  *ViewModelFactory
	WebFS       fs.FS

	Err404Template *template.Template
	IndexTemplate  *template.Template
	PostTemplate   *template.Template
}

func LoadDependencies(cfg Config) (*Dependencies, error) {
	deps := &Dependencies{}

	templateFS, err := fs.Sub(templateFiles, "templates")
	if err != nil {
		return nil, fmt.Errorf("sub embedded template dir: %w", err)
	}

	deps.WebFS, err = fs.Sub(webFiles, "web")
	if err != nil {
		return nil, fmt.Errorf("sub embedded web files: %w", err)
	}

	deps.ContentRoot, err = os.OpenRoot(ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
	}

	deps.ViewModels, err = NewModelViewFactory(cfg.Locales)
	if err != nil {
		return nil, fmt.Errorf("create view model factory: %v", err)
	}

	deps.Err404Template, err = template.ParseFS(templateFS, "*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load 404 template parts: %w", err)
	}
	_, err = deps.Err404Template.ParseFS(templateFS, "pages/404.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load 404 template: %w", err)
	}
	deps.Err404Template = deps.Err404Template.Lookup("layout.gohtml")

	deps.IndexTemplate, err = template.ParseFS(templateFS, "*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load index template parts: %w", err)
	}
	_, err = deps.IndexTemplate.ParseFS(templateFS, "pages/index.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load index template: %w", err)
	}
	deps.IndexTemplate = deps.IndexTemplate.Lookup("layout.gohtml")

	deps.PostTemplate, err = template.ParseFS(templateFS, "*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load post template parts: %w", err)
	}
	_, err = deps.PostTemplate.ParseFS(templateFS, "pages/post.gohtml")
	if err != nil {
		return nil, fmt.Errorf("load post template: %w", err)
	}
	deps.PostTemplate = deps.PostTemplate.Lookup("layout.gohtml")

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
		} else if argsErr, ok := ToArgsError(err); ok {
			PrintArgsErr(argsErr)
		} else {
			log.Fatalf("Error: load config: %v", err)
		}
		return
	}

	deps, err := LoadDependencies(cfg)
	if err != nil {
		log.Fatalf("Error: load deps: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.GetHead)
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))
	router.Get("/", indexHandler(deps, cfg.Locales))
	router.Handle("/_css/*", http.FileServerFS(deps.WebFS))
	router.Handle("/_fonts/*", http.FileServerFS(deps.WebFS))
	router.Handle("/_images/*", http.FileServerFS(deps.WebFS))
	router.Get("/*", contentHandler(deps))

	server := &http.Server{Addr: cfg.BindAddress, Handler: router}
	serverErr := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	log.Printf("Listening at %s\n", cfg.BindAddress)

	// Listen for CTRL+C or other interrupt.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server..")
	case err := <-serverErr:
		if err != nil {
			log.Printf("Error from server: %v", err)
		}
	}

	log.Println("Shutting down server...")
	stopCtx, stopStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopStop()
	err = server.Shutdown(stopCtx) // Shutdown server before closing deps.
	err = errors.Join(err, deps.Close())
	if err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Printf("Server stopped.")
}

func contentHandler(deps *Dependencies) http.HandlerFunc {
	versionURL := func(currentPath, target string) string {
		return web.VersionURL(deps.ContentRoot, currentPath, target)
	}
	parser := posts.NewParser(versionURL)

	return func(w http.ResponseWriter, r *http.Request) {
		response := web.NewResponse(w, r)
		filePath, _ := strings.CutPrefix(r.URL.Path, "/")

		// If no error, it is a regular file and not a post.
		if _, err := deps.ContentRoot.Stat(filePath); err == nil {
			response.File(deps.ContentRoot, filePath)
			return
		}

		if _, err := deps.ContentRoot.Stat(filePath + ".md"); err != nil {
			response.NotFound(deps.Err404Template, deps.ViewModels.Create(r, nil))
			return
		}

		post, err := parser.ParseFile(deps.ContentRoot.FS(), filePath+".md")
		if err != nil {
			response.ServerError(err)
			return
		}

		response.View(deps.PostTemplate, deps.ViewModels.Create(r, post))
	}
}

type SearchResults struct {
	Query   string
	Results []*posts.IndexedPost
}

func indexHandler(deps *Dependencies, locales []localization.Locale) http.HandlerFunc {
	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}

	indexes := make(map[string]posts.AutoIndex, len(locales))
	for _, loc := range locales {
		indexRoot := path.Join(ContentPath, loc.Tag)
		index, err := posts.NewAutoIndex(indexRoot)
		if err != nil {
			log.Fatalf("Error: load indexes: %v", err)
		}
		indexes[loc.Tag] = index
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqLocale := localization.LocaleFromContext(r.Context())
		reqQuery := r.URL.Query().Get("q")

		index := indexes[reqLocale.Tag]

		results := make([]*posts.IndexedPost, 0)
		for result := range index.Search(reqQuery) {
			if len(locales) != 0 {
				result.Path = "/" + path.Join(reqLocale.Tag, result.Path)
			}

			if result.Thumbnail != "" {
				root := deps.ContentRoot
				postPath := path.Dir(result.Path)
				thumbnail := web.VersionURL(root, postPath, result.Thumbnail)
				result.Thumbnail = path.Join(postPath, thumbnail)
			}

			results = append(results, result)
		}

		response := web.NewResponse(w, r)
		response.View(deps.IndexTemplate, deps.ViewModels.Create(r, SearchResults{
			Query:   reqQuery,
			Results: results,
		}))
	}
}
