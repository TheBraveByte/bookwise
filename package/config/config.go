package config

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/kataras/go-sessions/v3"
)

// CatalogueConfig : struct type which help hold field of reusable type in the
// application
type CatalogueConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Client      *http.Client
	Session     *sessions.Sessions
	Validate    *validator.Validate
}
