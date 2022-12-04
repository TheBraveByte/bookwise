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

	mux.Get("/view/books", c.AvailableBooks)
	mux.Post("/create/account", c.CreateAccount)
	mux.Post("/login/account", c.Login)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/user/search-book", c.SearchForBook)
		mux.Post("/user/pay/details", c.PurchaseBook)
		mux.Get("/user/pay/validate", c.ValidatePayment)
		mux.Get("/user/view/books", c.AvailableBooks)
		mux.Post("/user/view/library", c.ViewUserLibrary)
		mux.Get("/user/delete/book", c.DeleteUserBook)
		mux.Post("/user/search/book", c.SearchUserBook)
	})
	mux.Route("/add", func(mux chi.Router) {
		mux.Get("/new/book", c.AddBook)
	})

	return mux
}
