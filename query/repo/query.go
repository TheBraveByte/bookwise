package query

import (
	"context"
	"log"
	"time"

	"github.com/yusuf/p-catalogue/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (cr CatalogueDBRepo) CreateUserAccount(user model.User) (int, primitive.ObjectID, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: user.Email}}
	var user_data bson.M
	err := UserData(cr.DB, "user").FindOne(ctx, filter).Decode(&user_data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newUserID := primitive.NewObjectID()
			data := bson.D{
				{Key: "_id", Value: newUserID},
				{Key: "email", Value: user.Email},
				{Key: "first_name", Value: user.FirstName},
				{Key: "last_name", Value: user.LastName},
				{Key: "password", Value: user.Password},
				{Key: "catalogue", Value: user.Catalogue},
				{Key: "created_at", Value: user.CreatedAt},
			}
			_, err := UserData(cr.DB, "user").InsertOne(ctx, data)
			if err != nil {
				log.Fatal("cannot created database for user")
			}
			return 0, newUserID, nil
		}
		log.Fatal("this user does not exist")
	}
	var userID primitive.ObjectID
	for k, v := range user_data {
		if k == "_id" {
			switch id := v.(type) {
			case primitive.ObjectID:
				userID = id
			}
		}
	}
	return 1, userID, nil
}
