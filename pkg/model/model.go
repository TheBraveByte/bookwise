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
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Catalogue map[string]string  `json:"catalogue"`
}
