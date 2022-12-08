package config

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator"
)

// CatalogueConfig : struct type which help hold field of reusable type in the
// application
type CatalogueConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Session     *scs.SessionManager
	Client      *http.Client
	Validate    *validator.Validate
}
