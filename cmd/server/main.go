package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/subtributary/musings/internal/app"
)

func main() {
	config := loadConfig()
	if sane, err := config.IsSane(); !sane {
		log.Fatal(err)
	}

	server, err := app.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Dispose()

	fmt.Printf("Listening at %s\n", config.BindAddress)
	log.Fatal(server.ListenAndServe())
}

func loadConfig() *app.Config {
	config := app.NewConfig()

	config.BindAddress = os.Getenv("MUSINGS_BIND_ADDRESS")
	config.ContentPath = os.Getenv("MUSINGS_CONTENT_PATH")
	config.WebPath = os.Getenv("MUSINGS_WEB_PATH")

	flag.StringVar(&config.BindAddress, "web-endpoint", config.BindAddress, "Web endpoint to listen at")
	flag.Parse()

	return config
}
