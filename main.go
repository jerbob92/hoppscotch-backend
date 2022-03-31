package main

import (
	"log"

	"github.com/jerbob92/hoppscotch-backend/api"
	"github.com/jerbob92/hoppscotch-backend/config"
)

func init() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := api.StartAPI(); err != nil {
		log.Fatal(err)
	}
}
