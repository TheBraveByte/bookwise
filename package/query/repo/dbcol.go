package query

import "go.mongodb.org/mongo-driver/mongo"

func UserData(dbClient *mongo.Client, collection string) *mongo.Collection {
	var userData = dbClient.Database("p_catalogue").Collection(collection)
	return userData
}

func LibraryData(db *mongo.Client, collection string) *mongo.Collection {
	var library = db.Database("p_catalogue").Collection(collection)
	return library
}
