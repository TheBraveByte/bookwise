package query

import "go.mongodb.org/mongo-driver/mongo"

func UserData(dbClient *mongo.Client, collection string) *mongo.Collection {
	var userData = dbClient.Database("bookwise").Collection(collection)
	return userData
}

func LibraryData(db *mongo.Client, collection string) *mongo.Collection {
	var library = db.Database("bookwise").Collection(collection)
	return library
}
