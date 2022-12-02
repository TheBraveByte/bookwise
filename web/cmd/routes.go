package main

import (
	"github.com/yusuf/p-catalogue/modules/controller"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Route(catalog *controller.Catalogue) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(100 * time.Second))

	mux.Use(LoadAndSave)
	mux.Post("/create-account", catalog.CreateAccount)
	mux.Post("/login", catalog.Login)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/search-book", catalog.SearchBookTitle)

	})

	return mux
}
