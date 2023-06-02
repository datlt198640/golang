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
	// Setting up a Filter
	filter := bson.D{}
	// filter := bson.D{{"firstname", "Mili"}}
	// filter := bson.D{{"awards", bson.D{{"$gte", 10}}}}
	/*	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"awards", bson.D{{"$gt", 5}}}},
				bson.D{{"firstname", "Mili"}},
			}},
	}*/

	var results []*Actor

	// Finding multiple actors with a cursor
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem Actor
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor
	cur.Close(context.TODO())

	var i int
	for i = 0; i < len(results); i++ {
		fmt.Printf("Actor %v = %+v\n", i+1, *results[i])
	}

	// Close the connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("** Disconnected from MongoDB!")
	}

}
