package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/posts"
)

const (
	DataPath = "./data/"
)

func main() {
	var locale string
	flag.StringVar(&locale, "locale", "", "")
	flag.Usage = printUsage
	flag.Parse()
	if flag.NArg() != 1 {
		printUsage()
		os.Exit(1)
	}
	query := flag.Arg(0)

	index, err := posts.LoadIndex(DataPath, locale)
	if err != nil {
		log.Fatalf("Error loading index: %v", err)
	}

	count := 5
	for result := range index.Search(query) {
		fmt.Printf("%#v\n", result)

		count--
		if count == 0 {
			break
		}
	}
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Println("Musings post searcher.")
	fmt.Println()
	fmt.Printf("Usage: %s <query> [options]\n", program)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --locale <tag>  Set the locale for the index. [default: none]")
}
