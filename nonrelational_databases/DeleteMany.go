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
	filter := bson.D{}
	//filter := bson.D{{"firstname", "Mili"}}

	// Executing Delete operation
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v actors\n", deleteResult.DeletedCount)

	// Close the connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("** Disconnected from MongoDB!")
	}

}
