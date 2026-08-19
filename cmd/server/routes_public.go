package main

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/web"
)

func addPublicRoutes(router *chi.Mux, deps *Dependencies, cfg Config) {
	router.Get("/", indexHandler(deps, cfg.Locales))
	router.Handle("/_css/*", http.FileServerFS(deps.WebFS))
	router.Handle("/_fonts/*", http.FileServerFS(deps.WebFS))
	router.Handle("/_images/*", http.FileServerFS(deps.WebFS))
	router.Get("/*", contentHandler(deps))
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
			response.NotFound(
				web.WithView(deps.Templates["404.gohtml"]),
				web.WithData(deps.ViewModels.Create(r, nil)),
			)
			return
		}

		post, err := parser.ParseFile(deps.ContentRoot.FS(), filePath+".md")
		if err != nil {
			response.ServerError(err)
			return
		}

		response.Okay(
			web.WithView(deps.Templates["post.gohtml"]),
			web.WithData(deps.ViewModels.Create(r, post)),
		)
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
		indexRoot := path.Join(app.ContentPath, loc.Tag)
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
		vm := deps.ViewModels.Create(r, SearchResults{
			Query:   reqQuery,
			Results: results,
		})
		response.Okay(
			web.WithView(deps.Templates["index.gohtml"]),
			web.WithData(vm),
		)
	}
}
