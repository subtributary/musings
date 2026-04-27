package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"golang.org/x/text/language"
)

func main() {
	config := Config{}
	if err := config.LoadFromEnv(); err != nil {
		log.Fatalf("load from env: %v", err)
	}

	subcommand := ""
	if len(os.Args) >= 2 {
		subcommand = os.Args[1]
	}

	switch subcommand {
	case "rebuild":
		mainRebuild(config, os.Args[2:])
	case "rebuild-file":
		mainRebuildFile(config, os.Args[2:])
	case "search":
		mainSearch(config, os.Args[2:])
	default:
		log.Println("usage: indexer <subcommand>")
		log.Println("subcommand options are: rebuild, rebuild-file, search")
		os.Exit(1)
	}
}

func mainRebuild(config Config, args []string) {
}

func mainRebuildFile(config Config, args []string) {
	//
}

func mainSearch(config Config, args []string) {
	cmd := flag.NewFlagSet("rebuild", flag.ExitOnError)
	locale := cmd.String("locale", "", "locale to search")
	query := cmd.String("query", "", "search query")
	if err := cmd.Parse(args); err != nil {
		log.Fatalf("parse args: %v", err)
	}

	tag, err := language.Parse(*locale)
	if err != nil {
		log.Fatalf("parse locale: %v", err)
	}
	if !slices.Contains(config.Locales, tag) {
		log.Fatalf("locale not enabled: %s", *locale)
	}

	//
}
