package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/yusuf/p-catalogue/pkg/config"
	"github.com/yusuf/p-catalogue/pkg/controller"
	"github.com/yusuf/p-catalogue/pkg/database"
)

var app config.CatalogueConfig

func main() {
	mongoURI := os.Getenv("mongoURI")
	client := database.DatabaseConnection(mongoURI)
	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Panic(err)
		}
	}()

	catalog := controller.NewCatalogue(&app, client)
	controller.NewController(catalog)
	
	srv := &http.Server{Addr: ":8000", Handler: Route(&app)}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
