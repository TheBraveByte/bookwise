package main

import (
	"context"
	"encoding/gob"
	"github.com/yusuf/p-catalogue/api"
	"github.com/yusuf/p-catalogue/user/model"
	"github.com/yusuf/p-catalogue/user/userHandler"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/yusuf/p-catalogue/dependencies/config"
	"github.com/yusuf/p-catalogue/dependencies/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	app     config.CatalogueConfig
	session *scs.SessionManager
)

func main() {
	gob.Register(model.Data{})
	gob.Register(map[string]string{})
	gob.Register(primitive.NewObjectID())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("no environment variable file")
	}

	InfoLogger := log.New(os.Stdout, "p-catalogue info-logger", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, "p-catalogue error-logger : ", log.LstdFlags|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	//session.IdleTimeout = 60 * time.Minute
	session.Cookie.Persist = true
	session.Cookie.Secure = true
	session.Cookie.HttpOnly = true

	log.Println("Starting p-catalogue API application server ...............")

	uri := os.Getenv("mongodb_uri")
	client := database.DBConnection(uri)

	log.Println("....Application connected to the database.......")

	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Panic(err)
		}
	}()

	app.Session = session
	app.InfoLogger = InfoLogger
	app.ErrorLogger = ErrorLogger

	catalog := userHandler.NewCatalogue(&app, client)
	//user.NewController(catalog)

	libraryAPI := api.NewOpenLibraryAPI(&app, client)

	srv := &http.Server{Addr: ":8000", Handler: Route(catalog, libraryAPI)}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
