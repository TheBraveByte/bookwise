package model

type Book struct {
	Author   Author `json:"author"`
	Title    string `json:"title,omitempty"`
	Category string `json:"category"`
	ISBN     string `json:"ISBN"`
	Editor   string `json:"editor"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID        string            `bson:"_id" json:"_id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Catalogue map[string]string `json:"catalogue"`
}
