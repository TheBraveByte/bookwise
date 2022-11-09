package main

import (
	"github.com/yusuf/p-catalogue/api"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yusuf/p-catalogue/pkg/controller"
)

func Route(catalog *controller.Catalogue, api *api.OpenLibraryAPI) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(100 * time.Second))

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(LoadAndSave)
	mux.Post("/create-account", catalog.CreateAccount)
	mux.Post("/login", catalog.Login)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/search-book", api.SearchBookTitle)
	})

	return mux
}
