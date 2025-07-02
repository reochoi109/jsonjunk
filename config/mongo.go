package config

import (
	"context"
	logger "jsonjunk/pkg/logging"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var MongoClient *mongo.Client

func InitMongo() {
	logger.Log.Info("InitMongo [START]")
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
	MongoClient = client

	logger.Log.Info("InitMongo [END]",
		zap.String("host", "localhost"),
		zap.Int("port", 27017),
		zap.Bool("connected", true))
}
