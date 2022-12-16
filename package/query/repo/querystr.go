package query

import (
	"github.com/yusuf/bookwiseAPI/package/config"
	"github.com/yusuf/bookwiseAPI/package/query"
	"go.mongodb.org/mongo-driver/mongo"
)

type CatalogueDBRepo struct {
	App *config.CatalogueConfig
	DB  *mongo.Client
}

func NewCatalogueDBRepo(app *config.CatalogueConfig, db *mongo.Client) query.CatalogueRepo {
	return &CatalogueDBRepo{
		App: app,
		DB:  db,
	}
}
