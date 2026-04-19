package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelpAndExit()
	}

	switch os.Args[1] {
	case "index":
		mainIndex(os.Args)
	case "search":
		mainSearch(os.Args)
	default:
		printHelpAndExit()
	}
}

func mainIndex(args []string) {
	if len(args) != 2 && len(args) != 3 {
		printHelpAndExit()
	}

	config := IndexConfig{}
	if err := config.LoadFromEnv(); err != nil {
		log.Fatalf("invalid index config: %v", err)
	}
	config.LoadFromArgs(args[2:])

	// todo: index logic
}

func mainSearch(args []string) {
	if len(args) != 3 {
		printHelpAndExit()
	}

	config := SearchConfig{}
	if err := config.LoadFromEnv(); err != nil {
		log.Fatalf("invalid search config: %v", err)
	}
	config.LoadFromArgs(args[2:])

	// todo: search logic
}

func printHelpAndExit() {
	log.Print("Usage:")
	log.Print("  indexer index [path]")
	log.Print("  indexer search <query>")
	os.Exit(1)
}
