package database

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Database, error) {
	mongoUrl := os.Getenv("MONGODB_URL")
	if mongoUrl == "" {
		mongoUrl = "localhost:8000"
	}

	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database("membership"), err
}
