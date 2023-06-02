package main

import (
	"context"
	"fmt"
	"log"

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

	// Actor Details
	// james := Actor{"James", "Roger", 9}
	mark := Actor{"Mark", "Brown", 0}
	mili := Actor{"Mili", "Ford", 11}

	// Insert multiple documents
	actors := []interface{}{mark, mili}

	insertManyResult, err := collection.InsertMany(context.TODO(), actors)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple actors: ", insertManyResult.InsertedIDs)

	// Close the connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("** Disconnected from MongoDB!")
	}

}
