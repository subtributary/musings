package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/files"
	"github.com/subtributary/musings/internal/markdown"
	"github.com/subtributary/musings/internal/templates"
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

func loadConfig() *app.Config {
	config := app.NewConfig()

	config.BindAddress = os.Getenv("MUSINGS_BIND_ADDRESS")
	config.ContentPath = os.Getenv("MUSINGS_CONTENT_PATH")
	config.WebPath = os.Getenv("MUSINGS_WEB_PATH")

	flag.BoolVar(&config.EnableLiveTemplates, "live-templates", false, "Do not cache template files")
	flag.StringVar(&config.BindAddress, "web-endpoint", config.BindAddress, "Web endpoint to listen at")
	flag.Parse()

	return config
}

func loadServices(config *app.Config) (*app.Services, error) {
	services := app.Services{}

	if config.EnableLiveTemplates {
		services.TemplateProvider = templates.NewLiveTemplateProvider(config.GetTemplatesPath())
	} else {
		provider, err := templates.NewCachedTemplateProvider(config.GetTemplatesPath())
		if err != nil {
			return nil, fmt.Errorf("new template provider: %w", err)
		}
		services.TemplateProvider = provider
	}

	services.ContentStore = files.NewStore(config.ContentPath)
	services.MarkdownStore = markdown.NewStore(config.ContentPath)
	services.StaticStore = files.NewStore(config.GetStaticPath())

	return &services, nil
}
