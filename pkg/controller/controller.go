package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/yusuf/p-catalogue/pkg/config"
	"github.com/yusuf/p-catalogue/pkg/encrypt"
	"github.com/yusuf/p-catalogue/pkg/model"
	repo "github.com/yusuf/p-catalogue/query"
	query "github.com/yusuf/p-catalogue/query/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Validate = validator.New()

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
	var user model.User
	if err := rq.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	user.FirstName = rq.PostForm.Get("first_name")
	user.LastName = rq.PostForm.Get("last_name")
	user.Email = rq.PostForm.Get("email")
	user.Password, _ = encrypt.EncryptPassword(rq.PostForm.Get("password"))
	user.Catalogue = map[string]string{}
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().String())
	
	if err := Validate.Struct(user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			json.NewEncoder(wr).Encode(fmt.Sprintf("error %v %v", http.StatusBadRequest))
			return
		}
	}
	
	count, userID, _ := cg.CatDB.CreateUserAccount(user)
	ok := primitive.IsValidObjectID(userID.String())
	
	if count == 0 && ok {
		msg := model.ResponseMessage{
			StatusCode: http.StatusCreated,
			Message:    "Account created successfully",
		}
		json.NewEncoder(wr).Encode(msg)
		return
		
	} else {
		msg := model.ResponseMessage{
			StatusCode: http.StatusPermanentRedirect,
			Message:    "Existing account.Pls login",
		}
		json.NewEncoder(wr).Encode(msg)
		return
	}
}

func (cg *Catalogue) Login(wr http.ResponseWriter, rq *http.Request) {
	
}
