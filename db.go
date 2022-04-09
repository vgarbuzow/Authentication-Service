package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshToken struct {
	Guid    string
	Refresh []byte
	Time    time.Time
}

var collection *mongo.Collection
var ctx = context.TODO()

func initDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasker").Collection("tasks")
}

func createRefreshToken(refresh []byte, guid string) {
	token := RefreshToken{Refresh: refresh, Guid: guid, Time: time.Now().Add(60 * 24 * time.Hour)}
	insertResult, err := collection.InsertOne(context.TODO(), token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func readRefreshToken(guid string) {

}

/*func updateRefreshToken(refresh, guid string) {
	filter := bson.D{{"guid", guid}}
	update := bson.D{
		{"$set", bson.D{
			{"age", 1},
		}},
	}
}*/

func deleteRefreshToken(guid string) {

}
