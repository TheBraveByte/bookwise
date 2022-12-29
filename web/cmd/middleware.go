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
		// authToken, err := rq.Cookie("auth_token")
		// if err != nil {
		// 	log.Println(err)
		// }
		// fmt.Println()

		// if authToken.Value == "" {
		// 	log.Fatal("error no value is assigned to key in header")
		// 	return
		// }
		scs := session.Start(wr, rq)
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
		c, err := rq.Cookie("add_book")
		if err != nil {
			log.Println(err)
		}
		if c.Value == "" {
			log.Fatal("error no value is assigned to key in header")
			return
		}
		ctx := context.WithValue(rq.Context(), "purchase", c.Value)
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
