package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Library represents the main library for this application to store books
type Library struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id"`
	Book Book               `json:"book" bson:"book"`
}

// User represents the user for this application

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName   string             `json:"first_name" Usage:"required,alpha"`
	LastName    string             `json:"last_name" Usage:"required,alpha"`
	Email       string             `json:"email" Usage:"required,alphanumeric"`
	Password    string             `json:"password" Usage:"required,min=8,max=20"`
	UserLibrary []UserLibrary      `json:"user_library"`
	Token       string             `json:"token" Usage:"jwt"`
	RenewToken  string             `json:"renew_token" Usage:"jwt"`
	CreatedAt   time.Time          `json:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt   time.Time          `json:"updated_at" Usage:"datetime=2006-01-02"`
}

// UserLibrary : model to store purchased book by user

type UserLibrary struct {
	BookID     primitive.ObjectID `json:"_id" bson:"_id"`
	AuthorName []string           `json:"author_name"`
	Title      string             `json:"title"`
}

// Book sub-model for book API result
type Book struct {
	AuthorName   []string `json:"author_name"`
	Title        string   `json:"title"`
	PublishYear  int      `json:"first_publish_year"`
	Price        float64  `json:"price"`
	EditionCount int      `json:"edition_count"`
	Language     []string `json:"language"`
	Contributor  []string `json:"contributor"`
}

// Docs represent model for the API response when book is search by the user
type Docs struct {
	Docs []Book `json:"docs"`
}

// UserInfo represent model to store cookies data in sessions
type UserInfo struct {
	Email    string
	ID       primitive.ObjectID
	Password string
}

// PayLoad represent model struct for user payment details
type PayLoad struct {
	FirstName   string  `json:"first_name" Usage:"required,alpha"`
	LastName    string  `json:"last_name" Usage:"required,alpha"`
	Amount      float64 `json:"amount"`
	TxRef       string  `json:"tx_ref"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Currency    string  `json:"currency"`
	CardNo      string  `json:"card_no"`
	Cvv         string  `json:"cvv"`
	Pin         string  `json:"pin"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
}
