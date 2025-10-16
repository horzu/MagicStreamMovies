package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err:= godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	fmt.Println("MongoDB URI: ", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	err:= godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME environment variable not set")
	}
	fmt.Println("Database Name: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		log.Fatalf("Collection %s not found in database %s", collectionName, databaseName)
	}
	return collection
}