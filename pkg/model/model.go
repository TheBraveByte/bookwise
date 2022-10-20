package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID       string `json:"id"`
	Author   Author `json:"author"`
	Title    string `json:"title,omitempty"`
	Category string `json:"category"`
	ISBN     string `json:"ISBN"`
	Editor   string `json:"editor"`
	Price    string `json:"price"`
	Language string `json:"language"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName  string             `json:"first_name" Usage:"required,alpha"`
	LastName   string             `json:"last_name" Usage:"required,alpha"`
	Email      string             `json:"email" Usage:"required,alphanumeric"`
	Password   string             `json:"password" Usage:"required,min=8,max=20"`
	Catalogue  map[string]string  `json:"catalogue"`
	CreatedAt  string             `json:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt  string             `json:"updated_at" Usage:"datetime=2006-01-02"`
	Token      string             `json:"token" Usage:"jwt"`
	RenewToken string             `json:"renew_token" Usage:"jwt"`
}
