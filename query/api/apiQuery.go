package api

import (
	"context"
	"github.com/yusuf/p-catalogue/api"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

func (ap DBRepoAPI) AddBook(title string, bookData api.Book) (int64, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{"title", title}}
	count, err := BookData(ap.DB, "book").CountDocuments(ctx, filter)
	if err != nil {
		ap.App.ErrorLogger.Println(err)
		return 0, err
	}
	switch {
	case count == 0:
		newBook := bson.D{
			{"author_name", bookData.AuthorName},
			{"title", bookData.Title},
			{"first_publish_year", bookData.PublishYear},
			{"price", 99.55 * rand.Float64()},
			{"edition_count", bookData.EditionCount},
			{"language", bookData.Language},
			{"contributor", bookData.Contributor},
		}
		_, err := BookData(ap.DB, "book").InsertOne(ctx, newBook)
		if err != nil {
			ap.App.ErrorLogger.Fatalf("cannot created database for user %v", err)
		}
		return 0, nil
	case count >= 1:
		return 1, err

	}
	return 0, nil

}
