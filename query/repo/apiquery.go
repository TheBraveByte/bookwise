package query

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yusuf/p-catalogue/modules/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/rand"
	"time"
)

func (cr *CatalogueDBRepo) AddBook(title string, bookData model.Book) (int64, primitive.ObjectID, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{"title", title}}
	count, err := BookData(cr.DB, "book").CountDocuments(ctx, filter)
	if err != nil {
		cr.App.ErrorLogger.Println(err)
	}
	switch {
	case count == 0:
		newBook := bson.D{
			{"author_name", bookData.AuthorName},
			{"title", bookData.Title},
			{"first_publish_year", bookData.PublishYear},
			{"price", math.RoundToEven(99.55 * rand.Float64())},
			{"edition_count", bookData.EditionCount},
			{"language", bookData.Language},
			{"contributor", bookData.Contributor},
		}
		res, err := BookData(cr.DB, "book").InsertOne(ctx, newBook)
		if err != nil {
			cr.App.ErrorLogger.Fatalf("cannot created database for controller %v", err)
		}
		//var bookID primitive.ObjectID
		switch bookID := res.InsertedID.(type) {
		case primitive.ObjectID:
			return 0, bookID, nil
		}

	case count >= 1:
		return 1, [12]byte{}, nil
	}
	return 0, [12]byte{}, nil

}

func (cr *CatalogueDBRepo) GetSearchedBook(bookID primitive.ObjectID) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{"_id", bookID}}
	var res bson.M
	err := BookData(cr.DB, "book").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cr.App.ErrorLogger.Println(err)
			return nil, err
		}
		cr.App.ErrorLogger.Fatalln(err)
	}
	return res, nil
}
