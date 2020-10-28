package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func Initdb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		// stroe mongoURI in key.go
		mongoURI,
	))
	if err != nil { log.Fatal(err) }

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	database := client.Database("login_system_golang")
	usersCollection := database.Collection("users")
	fmt.Println(usersCollection)


	// test
	type User struct {
		User string
		Name string
		Tags    []string
	}
	user := User{
		User:  "The Polyglot Developer",
		Name: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
	}
	insertResult, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

}