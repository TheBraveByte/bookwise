package query

import (
	"github.com/yusuf/bookwiseAPI/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogueRepo interface {
	SendAvailableBooks() ([]primitive.M, error)

	CreateUserAccount(user *model.User) (int, primitive.ObjectID, error)
	VerifyUser(email, password, hashedPassword string) (bool, error)
	UpdateUserDetails(userID primitive.ObjectID, token, renewToken string) error

	UpdateUserBook(userID, bookId primitive.ObjectID, bookData primitive.M) error
	GetUserBooks(userID primitive.ObjectID) (primitive.A, error)
	FindBook(userID, bookId primitive.ObjectID) (primitive.M, error)
	DeleteBook(bookId primitive.ObjectID) error

	CheckLibrary(title string, bookData model.Book) (int64, primitive.ObjectID, error)
	GetBook(bookID primitive.ObjectID) (primitive.M, error)
}
