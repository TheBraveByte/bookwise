package query

import "go.mongodb.org/mongo-driver/mongo"

func UserData(dbClient *mongo.Client, _ string) *mongo.Collection {
	var userData = dbClient.Database("p_catalogue").Collection("user")
	return userData
}

func LibraryData(db *mongo.Client, collection string) *mongo.Collection {
	var library = db.Database("p_catalogue").Collection("library")
	return library
}
