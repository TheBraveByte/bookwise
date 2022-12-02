package config

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator"
)

type CatalogueConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Session     *scs.SessionManager
	Client      *http.Client
	Validate    *validator.Validate
}
