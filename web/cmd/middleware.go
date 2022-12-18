package main

import (
	"context"
	"github.com/yusuf/bookwiseAPI/package/token"
	"log"
	"net/http"
	"time"

	_"github.com/gorilla/securecookie"
	"github.com/kataras/go-sessions/v3"
)

// Authorization : middleware to authorize registered user to restricted routes using
// a unique generated token
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		authToken, err := rq.Cookie("auth_token")
		if err != nil {
			log.Println(err)
		}
		if authToken.Value == "" {
			log.Fatal("error no value is assigned to key in header")
			return
		}
		tokenClaims, err := token.ParseTokenString(authToken.Value)
		if err != nil {
			log.Fatalf("error %v", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(rq.Context(), "pass_token", tokenClaims)
		next.ServeHTTP(wr, rq.WithContext(ctx))

	})
}

// LoadAndSave : middleware for session to store cookies
//func LoadAndSave(next http.Handler) http.Handler {
//	return session.LoadAndSave(next)
//}

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


	


func CookieManager(cookieName string) *sessions.Sessions {
	scs := sessions.New(sessions.Config{
		Cookie:  cookieName,
		Expires: time.Hour * 24,
		// Encode:  secureCookie.Encode,
		// Decode:  secureCookie.Decode,
	})
	return scs
}