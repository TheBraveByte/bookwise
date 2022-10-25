package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/yusuf/p-catalogue/pkg/config"
	"github.com/yusuf/p-catalogue/pkg/controller"
	"github.com/yusuf/p-catalogue/pkg/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	app     config.CatalogueConfig
	session *scs.SessionManager
)

func main() {
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(primitive.NewObjectID())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("no environment variable file")
	}
	session = scs.New()
	session.Lifetime = 5 * time.Hour
	session.IdleTimeout = 60 * time.Minute
	session.Cookie.Path = "localhost"
	session.Cookie.Persist = true
	session.Cookie.Secure = true
	session.Cookie.HttpOnly = true

	log.Println("Starting p-catalogue API application server ...............")

	uri := os.Getenv("mongodb_uri")
	client := database.DatabaseConnection(uri)

	log.Println("....Application connected to the database.......")
	
	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Panic(err)
		}
	}()
	app.Session = session

	catalog := controller.NewCatalogue(&app, client)
	controller.NewController(catalog)

	srv := &http.Server{Addr: ":8000", Handler: Route(&app)}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
