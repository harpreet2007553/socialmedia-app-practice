package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var Collection_users *mongo.Collection
var Collection_posts *mongo.Collection
var Collection_comments *mongo.Collection


func ConnectDB(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	Collection_users = client.Database("SocialAppDB").Collection("users")
	Collection_posts = client.Database("SocialAppDB").Collection("posts")
	Collection_comments = client.Database("SocialAppDB").Collection("comments")

	
	fmt.Println("Successfully Connected To MONGODB!!")
}