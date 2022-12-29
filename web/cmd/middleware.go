package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/kataras/go-sessions/v3"
	"github.com/yusuf/bookwiseAPI/package/token"

	_ "github.com/gorilla/securecookie"
)

// Authorization : middleware to authorize registered user to restricted routes using
// a unique generated token
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {

		scs := app.Session.Start(wr, rq)
		authToken, ok := scs.Get("auth_token").(string)
		if !ok {
			log.Fatalf("%v token not available in session", http.StatusUnauthorized)
		}

		tokenClaims, err := token.ParseTokenString(authToken)
		if err != nil {
			log.Fatalf("error %v", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(rq.Context(), "pass_token", tokenClaims)
		next.ServeHTTP(wr, rq.WithContext(ctx))
	})
}

// LoadAndSave : middleware for session to store cookies

// AuthAddBook : middleware to authorize user to add books to their personal
// book collections
func AuthAddBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		scs := app.Session.Start(wr, rq)
		add, ok := scs.Get("add_book").(string)
		if !ok {
			log.Fatalf("%v token not available in session", http.StatusUnauthorized)
		}
		ctx := context.WithValue(rq.Context(), "purchased", add)
		next.ServeHTTP(wr, rq.WithContext(ctx))
	})
}

// CookieManager : Http session management setup to store cookies
func CookieManager(cookieName string) *sessions.Sessions {
	scs := sessions.New(sessions.Config{
		Cookie:  cookieName,
		Expires: time.Hour * 24,
		// Encode:  secureCookie.Encode,
		// Decode:  secureCookie.Decode,
	})
	return scs
}
