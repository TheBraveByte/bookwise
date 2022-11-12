package main

import (
	"github.com/yusuf/p-catalogue/modules/controller"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Route(catalog *controller.Catalogue) http.Handler {
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

	mux.Route("/api_repo", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/search-book", catalog.SearchBookTitle)
		//mux.Get("/get", catalog.APIServeUser)
	})

	return mux
}
