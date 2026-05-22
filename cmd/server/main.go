package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subtributary/musings/internal/localization"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(localization.LocalizedRoute(cfg.Locales))

	log.Printf("Listening at %s\n", cfg.BindAddress)
	log.Fatal(http.ListenAndServe(cfg.BindAddress, router))
}

//

/*
func loadServices(config Config) (services Services) {
	localization.InitTranslations()

	services.PostParser = posts.NewParser()

	services.PostIndexes = make(map[language.Tag]posts.Index)
	for _, tag := range config.Locales {
		locale := tag.String()
		if index, err := posts.LoadIndex(DataPath, locale); err != nil {
			log.Fatalf("could not load index: %v", err)
		} else {
			services.PostIndexes[tag] = index
		}
	}

	if config.EnableLiveTemplates {
		services.TemplateStore = templates.NewLiveStore(config.GetTemplatesPath(), config.Locales)
	} else {
		store, err := templates.NewCachedStore(config.GetTemplatesPath(), config.Locales)
		if err != nil {
			log.Fatalf("error loading templates: %v", err)
		}
		services.TemplateStore = store
	}

	return
}
*/
