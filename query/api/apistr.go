package api

import (
	"github.com/yusuf/p-catalogue/dependencies/config"
	"github.com/yusuf/p-catalogue/query"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBRepoAPI struct {
	App *config.CatalogueConfig
	DB  *mongo.Client
}

func NewDBRepoAPI(app *config.CatalogueConfig, db *mongo.Client) query.APIRepo {
	return &DBRepoAPI{
		App: app,
		DB:  db,
	}
}
