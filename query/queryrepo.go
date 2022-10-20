package repo

import (
	"github.com/yusuf/p-catalogue/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogueRepo interface {
	CreateUserAccount(user model.User) (int, primitive.ObjectID, error) 
}
