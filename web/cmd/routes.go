package main

import (
	"net/http"
	"time"

	"github.com/yusuf/bookwiseAPI/package/controller"

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

	mux.Get("/view/books", c.AvailableBooks)
	mux.Post("/create/account", c.CreateAccount)
	mux.Post("/login/account", c.Login)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(Authorization)
		mux.Post("/user/search-book", c.SearchForBook)
		mux.Post("/user/pay/details", c.PurchaseBook)
		mux.Get("/user/pay/validate", c.ValidatePayment)
		mux.Get("/user/view/library", c.AvailableBooks)
		mux.Get("/user/view/books", c.ViewUserLibrary)
		mux.Get("/user/delete/book/{id}", c.DeleteUserBook)
		mux.Get("/user/search/book/{id}", c.SearchUserBook)
	})
	mux.Route("/add", func(mux chi.Router) {
		mux.Use(AuthAddBook)
		mux.Get("/new/book", c.AddBook)
	})

	return mux
}
