package main

import (
	"net/http"
)
func Authorization(wr http.ResponseWriter, rq *http.Request){

}

func LoadAndSave(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
