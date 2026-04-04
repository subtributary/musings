package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/files"
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/store"
	"github.com/subtributary/musings/internal/templates"
	"golang.org/x/text/language"
)

func main() {
	config := loadConfig()
	if sane, err := config.IsSane(); !sane {
		log.Fatal(err)
	}

	services, err := loadServices(config)
	if err != nil {
		log.Fatal(err)
	}

	server, err := app.NewServer(services, config)
	if err != nil {
		log.Fatal(err)
	}

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
	config.Locales = tags

	return
}

func loadServices(config app.Config) (services app.Services, err error) {
	if config.EnableLiveTemplates {
		services.TemplateProvider = templates.NewLiveProvider(config.GetTemplatesPath(), config.Locales)
	} else {
		services.TemplateProvider, err = templates.NewCachedProvider(config.GetTemplatesPath(), config.Locales)
		if err != nil {
			err = fmt.Errorf("new template provider: %w", err)
			return
		}
	}

	services.ContentStore = files.NewStore(config.ContentPath, config.Locales)
	services.PostsStore = posts.NewStore(config.ContentPath, config.Locales)
	services.StaticStore = files.NewStore(config.GetStaticPath(), config.Locales)

	return
}
