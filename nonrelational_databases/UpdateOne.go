package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(">> Connected to MongoDB!")

	type Actor struct {
		FirstName string
		LastName  string
		Awards    int16
	}

	// Get a handle for your collection
	collection := client.Database("dvdstore").Collection("actordetails")

	// Setting up a Filter
	// filter := bson.D{}
	filter := bson.D{{"firstname", "Mili"}}
	// filter := bson.D{{"awards", bson.D{{"$gte", 10}}}}

	// Setting up a Update Value
	update := bson.D{
		{"$inc", bson.D{
			{"awards", 1},
		}},
	}

	// Executing Update operation
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v actors and updated %v actors.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Close the connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("** Disconnected from MongoDB!")
	}

}
