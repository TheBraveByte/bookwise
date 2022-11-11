package userHandler

import (
	"encoding/json"
	"fmt"
	"github.com/yusuf/p-catalogue/dependencies/token"
	"github.com/yusuf/p-catalogue/user/model"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/yusuf/p-catalogue/dependencies/config"
	"github.com/yusuf/p-catalogue/dependencies/encrypt"
	repo "github.com/yusuf/p-catalogue/query"
	query "github.com/yusuf/p-catalogue/query/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

var Validate = validator.New()

type Catalogue struct {
	App   *config.CatalogueConfig
	CatDB repo.CatalogueRepo
}

//var Catalog *Catalogue

func NewCatalogue(app *config.CatalogueConfig, db *mongo.Client) *Catalogue {
	return &Catalogue{
		App:   app,
		CatDB: query.NewCatalogueDBRepo(app, db),
	}
}

//
//func NewController(c *Catalogue) {
//	Catalog = c
//}

// CreateAccount : this function will help to create their account and have them store or add to
// database for future us
func (cg *Catalogue) CreateAccount(wr http.ResponseWriter, rq *http.Request) {
	var user model.User
	// Parse the posted details of the user
	if err := rq.ParseForm(); err != nil {
		cg.App.ErrorLogger.Fatalln(err)
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
			err := json.NewEncoder(wr).Encode(fmt.Sprintf("error %v", http.StatusBadRequest))
			if err != nil {
				return
			}
			return
		}
	}

	count, userID, _ := cg.CatDB.CreateUserAccount(user)

	// store user data in session as cookies
	data := model.Data{Email: user.Email, ID: userID, Password: user.Password}
	cg.App.Session.Put(rq.Context(), "data", data)

	switch {

	case count == 0:
		msg := model.ResponseMessage{
			StatusCode: http.StatusCreated,
			Message:    "account created successfully",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}

	case count == 1:
		msg := model.ResponseMessage{
			StatusCode: http.StatusPermanentRedirect,
			Message:    "existing account, pls login",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	}
}

// Login : this function will help to verify the user login details and
// also helps to generate authorization token for users
func (cg *Catalogue) Login(wr http.ResponseWriter, rq *http.Request) {

	if err := rq.ParseForm(); err != nil {
		cg.App.ErrorLogger.Fatalln(err)
	}
	email := rq.Form.Get("email")
	password := rq.Form.Get("password")

	data := cg.App.Session.Get(rq.Context(), "data").(model.Data)

	hashPassword := data.Password
	userID := fmt.Sprint(data.ID)

	switch {
	case email == data.Email:
		ok, _ := cg.CatDB.VerifyUser(email, password, hashPassword)
		if ok {
			generateToken, renewToken, err := token.GenerateToken(userID, email)
			if err != nil {
				log.Fatal(err)
				return
			}

			http.SetCookie(wr, &http.Cookie{Name: "auth_token", Value: generateToken, Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 1, 0)})

			_ = cg.CatDB.UpdateUserDetails(data.ID, generateToken, renewToken)

			msg := model.ResponseMessage{
				StatusCode: http.StatusOK,
				Message:    "You have login successfully",
			}
			err = json.NewEncoder(wr).Encode(msg)
			if err != nil {
				return
			}
		}
	default:
		msg := model.ResponseMessage{
			StatusCode: http.StatusUnauthorized,
			Message:    "Error: invalid login details",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	}
}
