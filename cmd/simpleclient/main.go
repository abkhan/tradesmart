package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// 1. Setup Configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect to MongoDB
	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect()

	// 3. Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
	fmt.Println("Successfully connected and pinged MongoDB!")

	// 4. Fetch movies (legacy logic preserved for this example)
	movies, err := findMoviesByRating(client.Client, "sample_mflix", "movies", "G")
	if err != nil {
		log.Fatalf("Error fetching movies: %v", err)
	}

	fmt.Printf("Found %d movies with rating 'G':\n", len(movies))
}

func findMoviesByRating(client *mongo.Client, dbName, collectionName, rating string) ([]bson.M, error) {
	coll := client.Database(dbName).Collection(collectionName)
	filter := bson.D{{Key: "rated", Value: rating}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []bson.M
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
