package main

import (
	"context"
	"github.com/yusuf/p-catalogue/modules/token"
	"log"
	"net/http"
)

// Authorization : middleware for authorization using the generated token
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
		ctx := context.WithValue(rq.Context(), "pass", tokenClaims)
		next.ServeHTTP(wr, rq.WithContext(ctx))
	})
}

// LoadAndSave : middleware for session to store cookies
func LoadAndSave(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
