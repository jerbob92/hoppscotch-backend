package main

import (
	"github.com/jerbob92/hoppscotch-backend/db"
	"github.com/jerbob92/hoppscotch-backend/models"
	"log"

	"github.com/jerbob92/hoppscotch-backend/api"
	"github.com/jerbob92/hoppscotch-backend/config"
	"github.com/jerbob92/hoppscotch-backend/fb"
)

func init() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	if err := models.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
	if err := fb.Initialize(); err != nil {
		log.Fatal(err)
	}
	if err := api.StartAPI(); err != nil {
		log.Fatal(err)
	}
}
