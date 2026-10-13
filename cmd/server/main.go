package main

import (
	"fmt"
	"log"

	"github.com/steviemul/gobi/internal/server"
)

func main() {

	const port = "8080"

	listener := &server.GobiListener{}

	if err := listener.Listen(port); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server listening on port %s, Ctrl+c to exit\n", port)

	if err := listener.Serve(); err != nil {
		log.Fatal(err)
	}
}
