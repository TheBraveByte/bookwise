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
				{Key: "user_library", Value: user.UserLibrary},
				{Key: "token", Value: ""},
				{Key: "renew_token", Value: ""},
				{Key: "created_at", Value: user.CreatedAt},
				{Key: "updated_at", Value: user.UpdatedAt},
			}
			_, err := UserData(cr.DB, "user").InsertOne(ctx, data)
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
	err := UserData(cr.DB, "user").FindOne(ctx, filter).Decode(&result)
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
		{"renew_token", renewToken}}}}

	_, err := UserData(cr.DB, "book").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CatalogueDBRepo) UpdateUserBook(userID primitive.ObjectID, book model.Library) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{"_id", userID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "user_library", Value: bson.D{
		{Key: "author_name", Value: book.Book.AuthorName},
		{Key: "title", Value: book.Book.Title},
		{Key: "book_id", Value: book.ID},
	}}}}}
	_, err := UserData(cr.DB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		cr.App.ErrorLogger.Fatalln(err)
	}
	return nil
}

func (cr *CatalogueDBRepo) GetUserBooks(userID primitive.ObjectID) ([]model.UserLibrary, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	var res model.User
	filter := bson.D{{"_id", userID}}
	err := UserData(cr.DB, "user").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		cr.App.ErrorLogger.Fatal(err)
	}
	return res.UserLibrary, nil
}
func (cr *CatalogueDBRepo) FindBook(userID primitive.ObjectID, title string) (model.UserLibrary, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "user_library.title", Value: title}, {Key: "_id", Value: userID}}
	opt := options.FindOne().SetProjection(bson.D{
		{Key: "title", Value: title},
		{Key: "user_library", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "title", Value: title}}}}},
	})
	var res model.User
	err := UserData(cr.DB, "user").FindOne(ctx, filter, opt).Decode(&res)
	if err != nil {
		cr.App.ErrorLogger.Panic(err)
	}
	return res.UserLibrary[0], nil

}

func (cr *CatalogueDBRepo) DeleteBook(bookID, userID primitive.ObjectID) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: userID}, {Key: "user_library.book_id", Value: bookID}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "user_library", Value: bson.D{{Key: "book_id", Value: bookID}}},
	}}}
	opt := options.Update().SetUpsert(false)

	_, err := UserData(cr.DB, "user").UpdateOne(ctx, filter, update, opt)
	if err != nil {
		cr.App.ErrorLogger.Fatalf("error when removing book from user library : %v ", err)
	}
	return nil
}
