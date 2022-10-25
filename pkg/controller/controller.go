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
	"github.com/yusuf/p-catalogue/token"

	repo "github.com/yusuf/p-catalogue/query"
	query "github.com/yusuf/p-catalogue/query/repo"
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

// CreateAccount : this function will help to create their account and have them store or add to 
// database for future us
func (cg *Catalogue) CreateAccount(wr http.ResponseWriter, rq *http.Request) {
	var user model.User
	// Parse the posted details of the user
	if err := rq.ParseForm(); err != nil {
		log.Fatal(err)
		return
	}

	// assigned parsed values to struct field and encrypt the input password
	user.FirstName = rq.PostForm.Get("first_name")
	user.LastName = rq.PostForm.Get("last_name")
	user.Email = rq.PostForm.Get("email")
	user.Password, _ = encrypt.EncryptPassword(rq.PostForm.Get("password"))
	user.Catalogue = map[string]string{}
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().String())

	// validate value with respect to struct tags
	if err := Validate.Struct(user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			json.NewEncoder(wr).Encode(fmt.Sprintf("error %v", http.StatusBadRequest))
			return
		}
	}

	count, userID, _ := cg.CatDB.CreateUserAccount(user)
	// ok := primitive.IsValidObjectID(userID.String())

	// store user data in session as cookies
	data := map[string]interface{}{"email": user.Email, "userID": userID, "password": user.Password}
	cg.App.Session.Put(rq.Context(), "data", data)

	switch {

	case count == 0:
		msg := model.ResponseMessage{
			StatusCode: http.StatusCreated,
			Message:    "account created successfully",
		}
		json.NewEncoder(wr).Encode(msg)

	case count == 1:
		msg := model.ResponseMessage{
			StatusCode: http.StatusPermanentRedirect,
			Message:    "existing account, pls login",
		}
		json.NewEncoder(wr).Encode(msg)
	}
}

// Login : this function will help to verify the user login details and 
// also helps to generate authorization token for users
func (cg *Catalogue) Login(wr http.ResponseWriter, rq *http.Request) {
	if err := rq.ParseForm(); err != nil {
		log.Fatal(err)
	}
	email := rq.PostForm.Get("email")
	password := rq.PostForm.Get("password")

	data := cg.App.Session.Get(rq.Context(), "data").(map[string]interface{})

	hashPassword := fmt.Sprintf("%s", data["password"])
	userID := fmt.Sprintf("%s", data["userID"])

	switch {
	case email == data["email"]:
		ok, _ := cg.CatDB.VerifyUser(email, password, hashPassword)
		if ok {
			token, _, err := token.GenerateToken(userID, email)
			if err != nil {
				log.Fatal(err)
				return
			}

			wr.Header().Set("Authorization", token)

			msg := model.ResponseMessage{
				StatusCode: http.StatusOK,
				Message:    "You have login successfully",
			}
			json.NewEncoder(wr).Encode(msg)
		}
	default:
		msg := model.ResponseMessage{
			StatusCode: http.StatusUnauthorized,
			Message:    "Error: invalid login details",
		}
		json.NewEncoder(wr).Encode(msg)
	}
}

func (cg *Catalogue) SearchBook(wr http.ResponseWriter, rq *http.Request) {

}
