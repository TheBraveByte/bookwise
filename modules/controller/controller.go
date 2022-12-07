package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/yusuf/p-catalogue/model"
	repo "github.com/yusuf/p-catalogue/modules/query"
	"github.com/yusuf/p-catalogue/modules/query/repo"
	"github.com/yusuf/p-catalogue/modules/token"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yusuf/p-catalogue/modules/config"
	"github.com/yusuf/p-catalogue/modules/encrypt"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:       3,
		DisableCompression: true,
		MaxConnsPerHost:    3,
		IdleConnTimeout:    100 * time.Second,
	},
}

type Catalogue struct {
	App    *config.CatalogueConfig
	CatDB  repo.CatalogueRepo
	Client *http.Client
}

func NewCatalogue(app *config.CatalogueConfig, db *mongo.Client) *Catalogue {
	return &Catalogue{
		App:    app,
		CatDB:  query.NewCatalogueDBRepo(app, db),
		Client: Client,
	}
}

// AvailableBooks : this method allows access for authorized and unauthorized users to views and
// see all available books in the library.
func (ct *Catalogue) AvailableBooks(wr http.ResponseWriter, rq *http.Request) {
	var library model.Library

	// validate value with respect to struct tags
	if err := ct.App.Validate.Struct(&library); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			err := json.NewEncoder(wr).Encode(fmt.Sprintf("error %v", http.StatusBadRequest))
			if err != nil {
				return
			}
			return
		}
	}
	books, err := ct.CatDB.SendAvailableBooks()
	if err != nil {
		ct.App.ErrorLogger.Fatalln(err)
	}

	msg := model.ResponseMessage{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("All available Books : \n %v ", books),
	}
	err = json.NewEncoder(wr).Encode(msg)
	if err != nil {
		return
	}
}

// CreateAccount : this methods will help to create their account and have them store or add to
// database for future usage.
func (ct *Catalogue) CreateAccount(wr http.ResponseWriter, rq *http.Request) {
	var user model.User

	// Parse the posted details of the controller
	if err := rq.ParseForm(); err != nil {
		ct.App.ErrorLogger.Fatalln(err)
		return
	}

	// assigned parsed values to struct field and encrypt the input password
	user.FirstName = rq.PostForm.Get("first_name")
	user.LastName = rq.PostForm.Get("last_name")
	user.Email = rq.PostForm.Get("email")
	user.Password, _ = encrypt.EncryptPassword(rq.PostForm.Get("password"))
	user.UserLibrary = []model.UserLibrary{}
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// validate value with respect to struct tags
	if err := ct.App.Validate.Struct(&user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			err := json.NewEncoder(wr).Encode(fmt.Sprintf("error %v", http.StatusBadRequest))
			if err != nil {
				return
			}
			return
		}
	}

	count, userID, _ := ct.CatDB.CreateUserAccount(user)

	// store controller data in session as cookies
	var data model.Data
	data.Email = user.Email
	data.ID = userID
	data.Password = user.Password

	ct.App.Session.Put(rq.Context(), "data", data)

	// check for new account or existing account
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

// Login : this method will help to verify the controller login details and also helps to generate
// authorization token for users.
func (ct *Catalogue) Login(wr http.ResponseWriter, rq *http.Request) {

	if err := rq.ParseForm(); err != nil {
		ct.App.ErrorLogger.Fatalln(err)
	}
	email := rq.Form.Get("email")
	password := rq.Form.Get("password")

	data := ct.App.Session.Get(rq.Context(), "data").(model.Data)

	hashPassword := data.Password
	userID := fmt.Sprint(data.ID)

	switch {
	case email == data.Email:
		ok, _ := ct.CatDB.VerifyUser(email, password, hashPassword)
		if ok {
			generateToken, renewToken, err := token.GenerateToken(userID, email)
			if err != nil {
				log.Fatal(err)
				return
			}

			http.SetCookie(wr, &http.Cookie{Name: "auth_token", Value: generateToken, Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 1, 0)})

			_ = ct.CatDB.UpdateUserDetails(data.ID, generateToken, renewToken)

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

// PurchaseBook : this will help users to process the payment procedure when the user provide
// their credit/debit card details.
func (ct *Catalogue) PurchaseBook(wr http.ResponseWriter, rq *http.Request) {

	// Parse the posted details of the controller
	if err := rq.ParseForm(); err != nil {
		ct.App.ErrorLogger.Fatalln(err)
		return
	}
	amount, _ := strconv.ParseFloat(rq.Form.Get("amount"), 64)
	payload := &model.PayLoad{
		FirstName:   rq.Form.Get("first_name"),
		LastName:    rq.Form.Get("last_name"),
		Amount:      amount,
		TxRef:       "MC-11001993",
		Email:       rq.Form.Get("email"),
		Phone:       rq.Form.Get("phone"),
		Currency:    rq.Form.Get("currency"),
		CardNo:      rq.Form.Get("card_no"),
		Cvv:         rq.Form.Get("cvv"),
		Pin:         rq.Form.Get("pin"),
		ExpiryMonth: rq.Form.Get("expiry_month"),
		ExpiryYear:  rq.Form.Get("expiry_year"),
	}

	// validate value with respect to struct tags
	if err := ct.App.Validate.Struct(payload); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			err := json.NewEncoder(wr).Encode(fmt.Sprintf("error %v", http.StatusBadRequest))
			if err != nil {
				return
			}
			return
		}
	}
	resp, _ := ct.Process(payload)

	ref := resp["data"].(map[string]interface{})["flwRef"].(string)

	ct.App.Session.Put(rq.Context(), "amount", amount)

	ct.App.Session.Put(rq.Context(), "ref", ref)

	err := json.NewEncoder(wr).Encode(&resp)
	if err != nil {
		return
	}

}

// ValidatePayment : this methods will help complete and verify the payment details of the user
// and other reference value needed.
func (ct *Catalogue) ValidatePayment(wr http.ResponseWriter, rq *http.Request) {

	ref := ct.App.Session.GetString(rq.Context(), "ref")
	if ref == "" {
		ct.App.ErrorLogger.Fatal("error no reference code for this transaction ")
	}
	amount := ct.App.Session.GetFloat(rq.Context(), "amount")
	if amount == 0.0 {
		ct.App.ErrorLogger.Fatal("error cannot make a payment of NGN 0.0 : pls enter a valid amount")
	}

	// Todo -> if this Application will be use in production
	// Todo ->  we need to get OTP input from the user during production of this API , but since we are using TEST MODE
	// Todo ->  we will be using a default OTP for now just for testing.

	resp, err := ct.Validate(ref, "1234")
	fmt.Println(resp)

	if err != nil {
		ct.App.ErrorLogger.Fatalf("invalid card details, cannot complete this payment process %v ", err)
	}

	respAmount := resp["data"].(map[string]interface{})["tx"].(map[string]interface{})["amount"].(float64)

	fmt.Printf(" amount is %v \n", respAmount)
	fmt.Printf(" status is %v \n", resp["status"])
	fmt.Printf(" message is %v \n", resp["message"])
	//fmt.Printf(" tX is %v \n", re)

	if resp["status"] == "success" && resp["message"] == "Charge Complete" && amount == respAmount {
		http.SetCookie(wr, &http.Cookie{Name: "add_book", Value: "success", Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 1, 0)})
		msg := model.ResponseMessage{
			StatusCode: http.StatusOK,
			Message:    "Payment Successful",
		}

		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	} else {
		msg := model.ResponseMessage{
			StatusCode: http.StatusInternalServerError,
			Message:    "Payment Not Successful ! Try again",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	}
}

// AddBook : this method will find/fetch the searched book using the bookID to extract book details
// from the database and make it available to the user; this is protected with a middleware
func (ct *Catalogue) AddBook(wr http.ResponseWriter, rq *http.Request) {

	bookID := ct.App.Session.Get(rq.Context(), "book_id").(primitive.ObjectID)
	fmt.Println(bookID)
	ok := primitive.IsValidObjectID(bookID.Hex())

	if !ok {
		ct.App.ErrorLogger.Println("book id is invalid")
	}

	book, err := ct.CatDB.GetBook(bookID)
	if err != nil {
		ct.App.ErrorLogger.Println("cannot find searched book")
	}

	bookData := book["book"].(primitive.M)
	bookId := book["_id"].(primitive.ObjectID)
	fmt.Println(book)
	//Add the book to the user Library
	data := ct.App.Session.Get(rq.Context(), "data").(model.Data)

	err = ct.CatDB.UpdateUserBook(data.ID, bookId, bookData)
	if err != nil {
		ct.App.ErrorLogger.Fatalln("error! cannot add book to user library")
	}

	err = json.NewEncoder(wr).Encode(&book)
	if err != nil {
		return
	}

}

// ViewUserLibrary : this method is to check out all the book collection that a particular user
// have bought
func (ct *Catalogue) ViewUserLibrary(wr http.ResponseWriter, rq *http.Request) {
	data := ct.App.Session.Get(rq.Context(), "data").(model.Data)
	userLibrary, err := ct.CatDB.GetUserBooks(data.ID)
	if err != nil {
		ct.App.ErrorLogger.Fatalln(err)
	}
	err = json.NewEncoder(wr).Encode(userLibrary)
	if err != nil {
		return
	}

}

// SearchUserBook : this method read data of a specific book title from the user library
func (ct *Catalogue) SearchUserBook(wr http.ResponseWriter, rq *http.Request) {
	//bookTitle:=chi.URLParam(rq, "title")
	title := rq.Form.Get("book_title")
	if title == "" {
		ct.App.ErrorLogger.Println("no input value to search")
		return
	}
	data := ct.App.Session.Get(rq.Context(), "data").(model.Data)
	book, err := ct.CatDB.FindBook(data.ID, title)
	if err != nil {
		ct.App.ErrorLogger.Fatal("error cannot find book")

	}
	err = json.NewEncoder(wr).Encode(book)
	if err != nil {
		return
	}

}

// DeleteUserBook : this will delete a specified book title in the user library
func (ct *Catalogue) DeleteUserBook(wr http.ResponseWriter, rq *http.Request) {
	title := rq.Form.Get("book_title")
	if title == "" {
		ct.App.ErrorLogger.Println("no input value to search")
		return
	}
	data := ct.App.Session.Get(rq.Context(), "data").(model.Data)
	book, err := ct.CatDB.FindBook(data.ID, title)
	if err != nil {
		ct.App.ErrorLogger.Fatal("error cannot find book")

	}
	bookID := book["_id"].(primitive.ObjectID)
	err = ct.CatDB.DeleteBook(bookID, data.ID)
	if err != nil {
		ct.App.ErrorLogger.Fatal("error cannot delete book")
	}
	msg := model.ResponseMessage{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Title : %v \nBook Deleted ", title),
	}
	err = json.NewEncoder(wr).Encode(msg)
	if err != nil {
		return
	}
}
