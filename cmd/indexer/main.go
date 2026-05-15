package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"golang.org/x/text/language"
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
	cfg, _ := configLoader.Load(args)

	if len(cfg.Locales) == 0 {
		log.Printf("no locales specified")
		return
	}

	root, err := os.OpenRoot(cfg.ContentPath)
	if err != nil {
		log.Fatalf("open root: %v", err)
	}
	defer func() { _ = root.Close() }()

	files := make(map[language.Tag][]string)
	dirQueue := []string{"/"}
	visited := make(map[string]struct{})
	for len(dirQueue) > 0 {
		dirName := dirQueue[0]

		dir, err := root.OpenRoot(dirName)
		if err != nil {
			log.Fatalf("open dir %q: %v", dirQueue[0], err)
		}

		result, err := localization.Scan(dir.FS())
		if err != nil {
			log.Fatalf("scan dir %q: %v", dirQueue[0], err)
		}

		for locale, entries := range result.GroupByTag(cfg.Locales) {
			for _, entry := range entries {
				path := filepath.Join(dirName, entry.Name())
				if entry.IsDir() {
					if _, ok := visited[path]; !ok {
						dirQueue = append(dirQueue, path)
						visited[path] = struct{}{}
					}
					continue
				}
				files[locale] = append(files[locale], path)
			}
		}
	}
}

/*
func rebuildDir(idx posts.Index, prefix string, dirFS fs.FS, locales []language.Tag) error {
	//
	return nil
}*/

func mainRebuildFile(args []string) {
	configLoader := config.NewLoader()
	configLoader.Presets(config.ContentPath, config.DataPath, config.Locales)
	configLoader.Required("target", "file to re-index")
	config, custom := configLoader.Load(args)
	target := custom["target"]

	log.Printf("config: %v", config)
	log.Printf("target: %v", target)
	log.Fatal("not implemented yet")
}

func mainSearch(args []string) {
	configLoader := config.NewLoader()
	configLoader.Presets(config.ContentPath, config.DataPath, config.Locales)
	configLoader.Required("query", "search query")
	config, custom := configLoader.Load(args)
	query := custom["query"]

	log.Printf("config: %v", config)
	log.Printf("query: %v", query)
	log.Fatal("not implemented yet")
}
