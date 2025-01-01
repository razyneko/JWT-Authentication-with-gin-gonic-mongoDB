package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBInstance() creates a mongo client
func DBinstance() *mongo.Client{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// creates a context with a deadline of 10 seconds from now. It ensures that any operation using this context will automatically be canceled if it exceeds the 10-second limit.

	defer cancel()
	// ensures that resources associated with the context are properly released when the function completes, avoiding potential resource leaks

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}

	// verify the connection
	err = client.Ping(ctx, nil)
	if err != nil{
		log.Fatal("Couldn't connect to MongoDB")
	}
	return client
}
// exportable var for mongo client
var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}
