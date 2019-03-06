package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DB = initDB()
)

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
