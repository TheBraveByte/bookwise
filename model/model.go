package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName  string             `json:"first_name" Usage:"required,alpha"`
	LastName   string             `json:"last_name" Usage:"required,alpha"`
	Email      string             `json:"email" Usage:"required,alphanumeric"`
	Password   string             `json:"password" Usage:"required,min=8,max=20"`
	Catalogue  map[string]string  `json:"catalogue"`
	CreatedAt  time.Time          `json:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt  time.Time          `json:"updated_at" Usage:"datetime=2006-01-02"`
	Token      string             `json:"token" Usage:"jwt"`
	RenewToken string             `json:"renew_token" Usage:"jwt"`
}

type Book struct {
	AuthorName   []string `json:"author_name"`
	Title        string   `json:"title"`
	PublishYear  int      `json:"first_publish_year"`
	Price        float64  `json:"price"`
	EditionCount int      `json:"edition_count"`
	Language     []string `json:"language"`
	Contributor  []string `json:"contributor"`
}

type Docs struct {
	Docs []Book `json:"docs"`
}

type ResponseMessage struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type Data struct {
	Email    string
	ID       primitive.ObjectID
	Password string
}
