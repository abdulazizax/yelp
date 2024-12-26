package main

import (
	"log"

	"github.com/abdulazizax/yelp/cmd/app"
	"github.com/abdulazizax/yelp/config"
)

func main() {
	// Load configuration and handle errors
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Run(cfg))
}
