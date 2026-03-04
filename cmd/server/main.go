package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/subtributary/musings/internal/app"
)

func main() {
	config := app.NewConfig()

	config.AssetsDir = os.Getenv("MUSINGS_ASSETS_DIR")
	config.ContentDir = os.Getenv("MUSINGS_CONTENT_DIR")
	config.WebEndpoint = os.Getenv("MUSINGS_WEB_ENDPOINT")

	flag.StringVar(&config.WebEndpoint, "web-endpoint", config.WebEndpoint, "Web endpoint to listen at")
	flag.Parse()

	if _, err := config.IsSane(); err != nil {
		log.Fatal(err)
	}

	server := app.NewServer(config)
	fmt.Printf("Listening at %s\n", config.WebEndpoint)
	log.Fatal(server.ListenAndServe())
}
