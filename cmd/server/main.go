package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server"
)

func main() {
	configPath := flag.String("config", "", "optional path to a YAML config file")
	serve := flag.Bool("serve", false, "serve the current directory as a browsable index (cannot be combined with -config)")
	flag.Parse()

	if *serve && *configPath != "" {
		fmt.Fprintln(os.Stderr, "-serve and -config are mutually exclusive")
		flag.Usage()
		os.Exit(2)
	}

	var cfg config.GobiConfig

	if *serve {
		cfg = config.IndexConfig()
	} else {
		var err error

		if cfg, err = config.BuildConfig(*configPath); err != nil {
			log.Fatal(err)
		}
	}

	listener := &server.GobiListener{}

	if err := listener.Listen(cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ctrl+c to exit")

	listener.Serve()
}
