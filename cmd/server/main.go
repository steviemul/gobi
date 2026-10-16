package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server"
)

func main() {
	configPath := flag.String("config", "", "optional path to a YAML config file")
	flag.Parse()

	config, err := config.BuildConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	listener := &server.GobiListener{}

	if err := listener.Listen(config); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ctrl+c to exit")

	listener.Serve()
}
