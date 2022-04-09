package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshToken struct {
	Guid    string
	Refresh string
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

	collection = client.Database("auth").Collection("refresh-tokens")
}

func createRefreshToken(refresh, guid string) {
	token := RefreshToken{Refresh: refresh, Guid: guid, Time: time.Now().Add(60 * 24 * time.Hour)}
	insertResult, err := collection.InsertOne(context.TODO(), token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func readRefreshToken(guid string) RefreshToken {
	filter := bson.D{{"guid", guid}}
	var result RefreshToken

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func updateRefreshToken(refresh, guid string) {
	filter := bson.D{{"guid", guid}}
	update := bson.D{
		{"$set", bson.D{
			{"refresh", refresh},
			{"time", time.Now().Add(60 * 24 * time.Hour)},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func deleteRefreshToken(guid string) {
	filter := bson.D{{"guid", guid}}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
