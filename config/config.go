package config

import "log"

type CatalogueConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}
