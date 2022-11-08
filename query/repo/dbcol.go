package query

import "go.mongodb.org/mongo-driver/mongo"

func UserData(dbClient *mongo.Client, _ string) *mongo.Collection {
	var userData = dbClient.Database("p_catalogue").Collection("user")
	return userData
}
