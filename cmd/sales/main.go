package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/mongodb"
	"mongotest/internal/parser"
)

const (
	dbName         = "tradesmart"
	collectionName = "sales"
)

func main() {
	filePath := flag.String("file", "", "Path to Sales CSV file")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide the sales data file path using the -file flag")
	}

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

	repo := mongodb.NewRepository(client.Database(dbName), collectionName)

	// 3. Process Input
	trades, err := parser.LoadTradesFromFile(*filePath)
	if err != nil {
		log.Fatalf("Error loading sales from file: %v", err)
	}

	if len(trades) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Ensure unique index on tracking
		if err := repo.CreateUniqueTrackingIndex(ctx); err != nil {
			log.Printf("Warning: Could not create unique index on tracking: %v", err)
		}

		// Perform bulk upsert
		res, err := repo.BulkUpsertByTracking(ctx, trades)
		if err != nil {
			log.Fatalf("Failed to bulk write sales: %v", err)
		}
		fmt.Printf("Successfully processed %d sales records from %s (Matched: %d, Upserted: %d, Modified: %d)\n",
			len(trades), *filePath, res.MatchedCount, res.UpsertedCount, res.ModifiedCount)
	}
}
