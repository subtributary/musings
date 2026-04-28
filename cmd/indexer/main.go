package main

import (
	"log"
	"os"

	"github.com/subtributary/musings/internal/config"
)

func main() {
	subcommand := ""
	if len(os.Args) >= 2 {
		subcommand = os.Args[1]
	}

	switch subcommand {
	case "rebuild":
		mainRebuild(os.Args[2:])
	case "rebuild-file":
		mainRebuildFile(os.Args[2:])
	case "search":
		mainSearch(os.Args[2:])
	default:
		log.Println("usage: indexer <subcommand>")
		log.Println("subcommand options are: rebuild, rebuild-file, search")
		os.Exit(1)
	}
}

func mainRebuild(args []string) {
	configLoader := config.NewLoader()
	configLoader.Presets(config.ContentPath, config.DataPath, config.Locales)
	config, _ := configLoader.Load(args)

	//
}

func mainRebuildFile(args []string) {
	configLoader := config.NewLoader()
	configLoader.Presets(config.ContentPath, config.DataPath, config.Locales)
	configLoader.Required("target", "file to re-index")
	config, custom := configLoader.Load(args)
	target := custom["target"]

	//
}

func mainSearch(args []string) {
	configLoader := NewConfigLoader()
	configLoader.Presets(ContentPath, DataPath, Locales)
	configLoader.Required("query", "search query")
	config, custom := configLoader.Load(args)
	query := custom["query"]

	//
}
