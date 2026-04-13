package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/templates"
	"golang.org/x/text/language"
)

func main() {
	config := loadConfig()
	if sane, err := config.IsSane(); !sane {
		log.Fatal(err)
	}

	services := loadServices(config)

	server := app.NewServer(services, config)

	fmt.Printf("Listening at %s\n", config.BindAddress)
	log.Fatal(server.ListenAndServe())
}

func loadConfig() (config app.Config) {
	config.BindAddress = os.Getenv("MUSINGS_BIND_ADDRESS")
	config.ContentPath = os.Getenv("MUSINGS_CONTENT_PATH")
	config.WebPath = os.Getenv("MUSINGS_WEB_PATH")
	locales := os.Getenv("MUSINGS_LOCALES")

	flag.BoolVar(&config.EnableLiveTemplates, "live-templates", false, "Do not cache template files")
	flag.StringVar(&config.BindAddress, "web-endpoint", config.BindAddress, "Web endpoint to listen at")
	flag.StringVar(&locales, "locales", locales, "Supported locales")
	flag.Parse()

	tags, _, err := language.ParseAcceptLanguage(locales)
	if err != nil {
		log.Fatalf("could not parse locales: %v", err)
	}
	if len(tags) == 0 {
		tags = []language.Tag{language.Und}
	}
	config.Locales = tags

	return
}

func loadServices(config app.Config) (services app.Services) {
	localization.InitTranslations()

	services.PostParser = posts.NewParser()

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
