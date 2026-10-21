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
	open := flag.Bool("open", false, "open the server in the default browser once it is listening (requires -serve)")
	flag.Parse()

	if *serve && *configPath != "" {
		usageError("-serve and -config are mutually exclusive")
	}

	if *open && !*serve {
		usageError("-open requires -serve")
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

	if *open {
		// Listen has already bound the port, so the browser can't beat the server.
		if err := openBrowser(browserURL(cfg)); err != nil {
			log.Printf("warning: could not open browser: %v", err)
		}
	}

	fmt.Println("Ctrl+c to exit")

	listener.Serve()
}

func usageError(message string) {
	fmt.Fprintln(os.Stderr, message)
	flag.Usage()
	os.Exit(2)
}
