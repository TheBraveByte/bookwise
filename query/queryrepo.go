package query

import (
	"github.com/yusuf/p-catalogue/modules/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogueRepo interface {
	CreateUserAccount(user model.User) (int, primitive.ObjectID, error)
	VerifyUser(email, password, hashedPassword string) (bool, error)
	UpdateUserDetails(userID primitive.ObjectID, token, renewToken string) error

	// AddBook controller Interacting with the book data
	AddBook(title string, bookData model.Book) (int64, primitive.ObjectID, error)
	//GetSearchedBook(bookID primitive.ObjectID) (primitive.M, error)
}
