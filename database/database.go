package database

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Database() *mongo.Client {
	godotenv.Load(".env")
	MONGODB_URL := os.Getenv("MONGODB_URL")

	mongoClient, _ := mongo.NewClient(options.Client().ApplyURI(MONGODB_URL))
	connection, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mongoClient.Connect(connection)
	return mongoClient
}

var MongoClient *mongo.Client = Database()

func Collection(client *mongo.Client, collectionName string) *mongo.Collection {
	var createCollection *mongo.Collection = client.Database("jwt-authentication").Collection(collectionName)

	return createCollection
}
