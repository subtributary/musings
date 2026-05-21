package main

func main() {
	/*
		var liveTemplates bool
		var address string
		flag.BoolVar(&liveTemplates, "live-templates", false, "")
		flag.StringVar(&address, "bind", DefaultBindAddress, "")
		flag.Usage = printUsage
		flag.Parse()
		if flag.NArg() > 1 {
			printUsage()
			os.Exit(1)
		}

		services := loadServices(config)

		server := NewServer(services, config)

		fmt.Printf("Listening at %s\n", config.BindAddress)
		log.Fatal(server.ListenAndServe())*/
}

/*
func loadServices(config Config) (services Services) {
	localization.InitTranslations()

	services.PostParser = posts.NewParser()

	services.PostIndexes = make(map[language.Tag]posts.Index)
	for _, tag := range config.Locales {
		locale := tag.String()
		if index, err := posts.LoadIndex(DataPath, locale); err != nil {
			log.Fatalf("could not load index: %v", err)
		} else {
			services.PostIndexes[tag] = index
		}
	}

	if config.EnableLiveTemplates {
		services.TemplateStore = templates.NewLiveStore(config.GetTemplatesPath(), config.Locales)
	} else {
		store, err := templates.NewCachedStore(config.GetTemplatesPath(), config.Locales)
		if err != nil {
			log.Fatalf("error loading templates: %v", err)
		}
		services.TemplateStore = store
	}

	return
}
*/
