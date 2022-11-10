package api

import "go.mongodb.org/mongo-driver/mongo"

func BookData(db *mongo.Client, collection string) *mongo.Collection {
	var bookCollection = db.Database("p_catalogue").Collection("book")
	return bookCollection
}
