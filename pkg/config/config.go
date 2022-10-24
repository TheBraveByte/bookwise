package config

import (
	"log"

	"github.com/alexedwards/scs/v2"
)

type CatalogueConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Session     *scs.SessionManager
}
