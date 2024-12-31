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
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database("membership"), err
}
