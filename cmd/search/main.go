package main

import (
	"fmt"
	"log"

	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/posts"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	index, err := posts.LoadIndex(config.DataPath, cfg.Locale)
	if err != nil {
		log.Fatalf("Error loading index: %v", err)
	}

	count := 5
	for result := range index.Search(cfg.Query) {
		fmt.Printf("%#v\n", result)

		count--
		if count == 0 {
			break
		}
	}
}
