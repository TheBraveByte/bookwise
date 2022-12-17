package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator"
	"github.com/yusuf/bookwiseAPI/database"
	"github.com/yusuf/bookwiseAPI/model"
	"github.com/yusuf/bookwiseAPI/package/controller"

	"github.com/joho/godotenv"
	"github.com/yusuf/bookwiseAPI/package/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validate *validator.Validate
	app      config.CatalogueConfig
	session  *scs.SessionManager
)

func main() {
	gob.Register(model.UserInfo{})
	gob.Register(model.PayLoad{})
	gob.Register(model.User{})
	gob.Register(model.Book{})
	gob.Register(model.Docs{})
	gob.Register(model.UserLibrary{})
	gob.Register(map[string]string{})
	gob.Register(primitive.NewObjectID())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("no environment variable file")
	}

	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	validate = validator.New()
	app.Validate = validate

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Domain = "localhost"
	// session.Cookie.Path = "/"
	// session.Cookie.HttpOnly = true

	app.Session = session

	app.InfoLogger = InfoLogger
	app.ErrorLogger = ErrorLogger

	log.Println("..........  Starting Bookwise API application server  ..........")
	uri := os.Getenv("MONGODB_URI")

	client := database.Connection(uri)

	log.Println("..........  Application connected to the database  ..........")

	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Panic(err)
		}
	}()

	catalog := controller.NewCatalogue(&app, client)

	srv := &http.Server{Addr: ":8000", Handler: Route(catalog)}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
