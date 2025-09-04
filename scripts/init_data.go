
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	ID    int32  `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Year  string `json:"year" bson:"year"`
}

func main() {
	// Get MongoDB connection string from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password@mongodb:27017/movies_db?authSource=admin"
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "movies_db"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")

	// Read movies.json file
	data, err := os.ReadFile("/app/movies.json")
	if err != nil {
		log.Fatalf("Failed to read movies.json: %v", err)
	}

	var movies []Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		log.Fatalf("Failed to parse movies.json: %v", err)
	}

	fmt.Printf("Loaded %d movies from JSON file\n", len(movies))

	// Get database and collection
	db := client.Database(databaseName)
	collection := db.Collection("movies")

	// Check if data already exists
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to count existing documents: %v", err)
	}

	if count > 0 {
		fmt.Printf("Database already contains %d movies. Skipping initialization.\n", count)
		return
	}

	// Convert to interface slice for bulk insert
	docs := make([]interface{}, len(movies))
	for i, movie := range movies {
		docs[i] = movie
	}

	// Insert movies
	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to insert movies: %v", err)
	}

	fmt.Printf("Successfully inserted %d movies into the database!\n", len(result.InsertedIDs))

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "title", Value: "text"}, {Key: "year", Value: 1}},
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	} else {
		fmt.Println("Successfully created database indexes!")
	}

	fmt.Println("Database initialization completed successfully!")
}

