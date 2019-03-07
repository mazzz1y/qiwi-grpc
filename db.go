package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DB = initDB()
)

// Create unique index if number of rows = 1
func SetUniqueIndex(collection *mongo.Collection, rowName string) error {
	allDocCount, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	if allDocCount == 1 {
		_, err = collection.Indexes().CreateOne(
			context.TODO(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{rowName, bsonx.Int32(1)}},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			return err
		}
	}
	return nil
}


func initDB() *mongo.Database {
	uri := "mongodb://" + Config.MongoHost + ":" + Config.MongoPort
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if client.Ping(ctx, readpref.Primary()) != nil {
		log.Fatal("MongoDB: connect failed")
	}
	log.Println("MongoDB: connected: " + uri)
	return client.Database(Config.MongoDB)
}
