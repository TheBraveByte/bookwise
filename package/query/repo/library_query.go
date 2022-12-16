package query

import (
	"context"
	"github.com/yusuf/bookwiseAPI/model"
	"go.mongodb.org/mongo-driver/mongo"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/rand"
	"time"
)

func (cr *CatalogueDBRepo) SendAvailableBooks() ([]primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{}}
	var res []bson.M
	cursor, err := LibraryData(cr.DB, "library").Find(ctx, filter)
	if err != nil {
		cr.App.ErrorLogger.Fatalf("error in library collection : %v ", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// CheckLibrary : this query method will check if a particular book searched is in the
// library collection and if it is not found in the library insert/add the details of the
// book in the library
func (cr *CatalogueDBRepo) CheckLibrary(title string, bookData model.Book) (int64, primitive.ObjectID, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "book.title", Value: title}}
	var res model.Library
	err := LibraryData(cr.DB, "library").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			bookID := primitive.NewObjectID()
			newBook := bson.D{
				{Key: "_id", Value: bookID},
				{Key: "book", Value: bson.D{
					{Key: "author_name", Value: bookData.AuthorName},
					{Key: "title", Value: bookData.Title},
					{Key: "first_publish_year", Value: bookData.PublishYear},
					{Key: "price", Value: math.RoundToEven(99.55 * (rand.Float64() + 5))},
					{Key: "edition_count", Value: bookData.EditionCount},
					{Key: "language", Value: bookData.Language},
					{Key: "contributor", Value: bookData.Contributor},
				}}}
			_, err := LibraryData(cr.DB, "library").InsertOne(ctx, newBook)
			if err != nil {
				cr.App.ErrorLogger.Fatalf("cannot created database for controller %v", err)
			}

			return 0, bookID, nil
		}
		log.Fatal(err)

	}
	return 1, res.ID, nil

}

func (cr *CatalogueDBRepo) GetBook(bookID primitive.ObjectID) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{"_id", bookID}}
	var res bson.M
	err := LibraryData(cr.DB, "library").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		cr.App.ErrorLogger.Fatalln(err)
	}
	return res, nil
}
