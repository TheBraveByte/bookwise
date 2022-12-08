package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/yusuf/p-catalogue/model"
	repo "github.com/yusuf/p-catalogue/package/query"
	"github.com/yusuf/p-catalogue/package/query/repo"
	"github.com/yusuf/p-catalogue/package/token"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yusuf/p-catalogue/package/config"
	"github.com/yusuf/p-catalogue/package/encrypt"
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

	numberOfBooks := len(books)
	var authorName []primitive.A
	for _, m := range books {
		author := m["book"].(primitive.M)["author_name"].(primitive.A)
		authorName = append(authorName, author)
	}

	fmt.Println(authorName)

	numberOfAuthors := len(authorName)

	stat := map[string]interface{}{
		"status_code":       http.StatusOK,
		"message":           "All available books in library",
		"number_of_authors": numberOfAuthors,
		"number_of_books":   numberOfBooks,
		"data":              books,
	}

	err = json.NewEncoder(wr).Encode(stat)
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
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

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

	// store controller userInfo in session as cookies
	userInfo := model.UserInfo{
		Email:    user.Email,
		ID:       userID,
		Password: user.Password,
	}

	ct.App.Session.Put(rq.Context(), "userInfo", userInfo)

	// check for new account or existing account
	switch {
	case count == 0:
		msg := map[string]interface{}{
			"status_code": http.StatusCreated,
			"message":     "Account Registered Successfully",
		}

		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}

	case count == 1:
		msg := map[string]interface{}{
			"status_code": http.StatusPermanentRedirect,
			"message":     "Existing Account; Please Login",
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

	userInfo := ct.App.Session.Get(rq.Context(), "userInfo").(model.UserInfo)

	hashPassword := userInfo.Password
	userID := fmt.Sprint(userInfo.ID)

	switch {
	case email == userInfo.Email:
		ok, _ := ct.CatDB.VerifyUser(email, password, hashPassword)
		if ok {
			generateToken, renewToken, err := token.GenerateToken(userID, email)
			if err != nil {
				log.Fatal(err)
				return
			}

			http.SetCookie(wr, &http.Cookie{Name: "auth_token", Value: generateToken, Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 1, 0)})

			err = ct.CatDB.UpdateUserDetails(userInfo.ID, generateToken, renewToken)
			if err != nil {
				ct.App.ErrorLogger.Fatalln("error cannot update token ")
				return
			}

			msg := map[string]interface{}{
				"status_code": http.StatusOK,
				"message":     "Successfully Logged-in",
			}
			err = json.NewEncoder(wr).Encode(msg)
			if err != nil {
				return
			}
		}
	default:
		msg := map[string]interface{}{
			"status_code": http.StatusUnauthorized,
			"message":     "Error !!! : Invalid Login Details",
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

	msg := map[string]interface{}{
		"status_code": http.StatusAccepted,
		"message":     "All Payment Details Confirmed",
		"response":    resp,
	}
	err := json.NewEncoder(wr).Encode(msg)
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
	if err != nil {
		ct.App.ErrorLogger.Fatalf("invalid card details, cannot complete this payment process %v ", err)
	}

	respAmount := resp["data"].(map[string]interface{})["tx"].(map[string]interface{})["amount"].(float64)

	fmt.Printf(" amount is %v \n", respAmount)
	fmt.Printf(" status is %v \n", resp["status"])
	fmt.Printf(" message is %v \n", resp["message"])

	if resp["status"] == "success" && resp["message"] == "Charge Complete" && amount == respAmount {
		http.SetCookie(wr, &http.Cookie{Name: "add_book", Value: "success", Path: "/", Domain: "localhost", Expires: time.Now().Add(60 * time.Second)})
		msg := map[string]interface{}{
			"status_code": http.StatusOK,
			"message":     "Payment For Book Successful",
			"response":    resp,
		}

		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	} else {
		msg := map[string]interface{}{
			"status_code": http.StatusInternalServerError,
			"message":     "Payment Not Successful ! Try again",
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

	//Add the book to the user Library
	userInfo := ct.App.Session.Get(rq.Context(), "userInfo").(model.UserInfo)

	err = ct.CatDB.UpdateUserBook(userInfo.ID, bookId, bookData)
	if err != nil {
		ct.App.ErrorLogger.Fatalln("error! cannot add book to user library")
	}
	msg := map[string]interface{}{
		"status_code": http.StatusOK,
		"message":     "New Book Added To Library",
		"data":        book,
	}
	err = json.NewEncoder(wr).Encode(msg)
	if err != nil {
		return
	}

}

// ViewUserLibrary : this method is to check out all the book collection that a particular user
// have bought
func (ct *Catalogue) ViewUserLibrary(wr http.ResponseWriter, rq *http.Request) {
	userInfo := ct.App.Session.Get(rq.Context(), "userInfo").(model.UserInfo)
	userLibrary, err := ct.CatDB.GetUserBooks(userInfo.ID)
	if err != nil {
		ct.App.ErrorLogger.Fatalln(err)
	}
	msg := map[string]interface{}{
		"status_code": http.StatusOK,
		"message":     "User Book Collections",
		"stat":        len(userLibrary),
		"data":        userLibrary,
	}
	err = json.NewEncoder(wr).Encode(msg)
	if err != nil {
		return
	}

}

// SearchUserBook : this method read data of a specific book title from the user library
func (ct *Catalogue) SearchUserBook(wr http.ResponseWriter, rq *http.Request) {
	bookID, err := primitive.ObjectIDFromHex(chi.URLParam(rq, "id"))
	if err != nil {
		ct.App.ErrorLogger.Fatalln("invalid id parameter")
		return
	}
	userInfo := ct.App.Session.Get(rq.Context(), "userInfo").(model.UserInfo)
	book, err := ct.CatDB.FindBook(userInfo.ID, bookID)
	if err != nil {
		ct.App.ErrorLogger.Fatal("error cannot find book")

	}
	if len(book) == 0 {
		msg := map[string]interface{}{
			"status_code": http.StatusOK,
			"message":     "Book Not Found",
			"data":        book,
		}

		err = json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}

	}
	if len(book) >= 1 {
		msg := map[string]interface{}{
			"status_code": http.StatusOK,
			"message":     "Book is Found",
			"data":        book,
		}

		err = json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}

	}

}

// DeleteUserBook : this will delete a specified book title in the user library
func (ct *Catalogue) DeleteUserBook(wr http.ResponseWriter, rq *http.Request) {
	bookID, err := primitive.ObjectIDFromHex(chi.URLParam(rq, "id"))
	if err != nil {
		ct.App.ErrorLogger.Fatalln("invalid book id")
		return
	}

	err = ct.CatDB.DeleteBook(bookID)
	if err != nil {
		ct.App.ErrorLogger.Fatal("error cannot delete book")
	}
	msg := map[string]interface{}{
		"status_code": http.StatusOK,
		"message":     fmt.Sprintf("Book with an ID: %v Deleted", bookID),
	}
	err = json.NewEncoder(wr).Encode(msg)
	if err != nil {
		return
	}
}
