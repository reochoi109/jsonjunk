package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://admin:admin123@localhost:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("mongodb connection failed , ", err)
		return
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongodb ping failed , ", err)
	}

	fmt.Println("connection success")
	MongoClient = client
}
