package controller

import (
	"net/http"

	"github.com/yusuf/p-catalogue/pkg/config"
	repo "github.com/yusuf/p-catalogue/query"
	query "github.com/yusuf/p-catalogue/query/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Catalogue struct {
	App   *config.CatalogueConfig
	CatDB repo.CatalogueRepo
}

var Catalog *Catalogue

func NewCatalogue(app *config.CatalogueConfig, db *mongo.Client) *Catalogue {
	return &Catalogue{
		App:   app,
		CatDB: query.NewCatalogueDBRepo(app, db),
	}
}

func NewController(c *Catalogue) {
	Catalog = c
}

func (cg *Catalogue) CreateAccount(wr http.ResponseWriter, rq *http.Request) {
	
}

func (cg *Catalogue) Login(wr http.ResponseWriter, rq *http.Request) {
}
