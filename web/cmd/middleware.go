package main

import (
	"context"
	"log"
	"net/http"

	"github.com/yusuf/p-catalogue/token"
)
//Authorization : middleware for authorization using the generated token
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		authToken := rq.Header.Get("Authorization")
		if authToken == "" {
			log.Fatal("error no value is assigned to key in header")
			return
		}
		tokenClaims, err := token.ParseTokenString(authToken)
		if err != nil {
			log.Fatalf("error %v", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(rq.Context(),"pass", tokenClaims)
		next.ServeHTTP(wr, rq.WithContext(ctx))
	})
}

//LoadAndSave : middleware for session to store cookies
func LoadAndSave(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
