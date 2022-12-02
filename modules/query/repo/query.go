package query

import (
	"context"
	"github.com/yusuf/p-catalogue/model"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/yusuf/p-catalogue/modules/encrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (cr *CatalogueDBRepo) CreateUserAccount(user model.User) (int, primitive.ObjectID, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: user.Email}}
	var userData bson.M
	err := UserData(cr.DB, "controller").FindOne(ctx, filter).Decode(&userData)
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
			_, err := UserData(cr.DB, "controller").InsertOne(ctx, data)
			if err != nil {
				cr.App.ErrorLogger.Fatalf("cannot created database for controller %v", err)
			}
			return 0, newUserID, nil
		}
		cr.App.ErrorLogger.Fatalf("this controller does not exist, %v", err)
	}
	var userID primitive.ObjectID
	for k, v := range userData {
		if k == "_id" {
			switch id := v.(type) {
			case primitive.ObjectID:
				userID = id
			}
		}
	}
	return 1, userID, nil
}

func (cr *CatalogueDBRepo) VerifyUser(email, password, encryptPassword string) (bool, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: email}}
	var result bson.M
	err := UserData(cr.DB, "controller").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, err
		}
		return false, err
	}

	ok, err := encrypt.VerifyEncryptPassword(password, encryptPassword)
	if err != nil {
		cr.App.ErrorLogger.Println(err)
	}
	return ok, nil
}

func (cr *CatalogueDBRepo) UpdateUserDetails(userID primitive.ObjectID, token, renewToken string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$set", bson.D{
		{"token", token},
		{"renew_token", renewToken},
	}}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var newUpdate bson.M
	err := UserData(cr.DB, "book").FindOneAndUpdate(ctx, filter, update, opts).Decode(&newUpdate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		cr.App.ErrorLogger.Fatal(err)
	}
	return nil

}
