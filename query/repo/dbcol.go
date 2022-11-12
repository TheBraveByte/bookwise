package query

import "go.mongodb.org/mongo-driver/mongo"

func UserData(dbClient *mongo.Client, _ string) *mongo.Collection {
	var userData = dbClient.Database("p_catalogue").Collection("controller")
	return userData
}

func BookData(db *mongo.Client, collection string) *mongo.Collection {
	var bookCollection = db.Database("p_catalogue").Collection("book")
	return bookCollection
}
