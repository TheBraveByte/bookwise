package main

import (
	"github.com/yusuf/p-catalogue/modules/controller"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Route(c *controller.Catalogue) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(100 * time.Second))

	//middleware for data store in session as cookies
	mux.Use(LoadAndSave)

	//endpoint
	mux.Get("/view-books", c.AvailableBooks)
	mux.Post("/create/account", c.CreateAccount)
	mux.Post("/login/account", c.Login)
	mux.Post("/pay/book", c.PurchaseBook)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/search-book", c.SearchBookTitle)
		mux.Post("/pay/details", c.PurchaseBook)
		mux.Get("/pay/validate", c.ValidatePayment)

		mux.Use(AuthAddBook)
		mux.Get("/add/book", c.AddBook)
	})

	return mux
}
