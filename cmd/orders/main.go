package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/models"
	"mongotest/internal/mongodb"
	"mongotest/internal/parser"
)

const (
	dbName         = "tradesmart"
	collectionName = "orders"
)

func main() {
	// CLI Flags
	product := flag.String("product", "", "Product Name")
	qty := flag.Int("qty", 0, "Quantity")
	price := flag.Float64("price", 0.0, "Sale Price")
	orderType := flag.String("type", "Sale", "Order Type")
	sellerID := flag.String("id", "", "Seller Order ID")
	filePath := flag.String("file", "", "Path to CSV or Excel file")
	flag.Parse()

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
	var trades []models.Trade

	if *filePath != "" {
		trades, err = parser.LoadTradesFromFile(*filePath)
		if err != nil {
			log.Fatalf("Error loading trades from file: %v", err)
		}
	} else if *product != "" && *qty > 0 {
		trade := models.Trade{
			ProductName:   *product,
			Qty:           *qty,
			SalePrice1:    *price,
			OrderType:     *orderType,
			SellerOrderID: *sellerID,
			OrderDate:     time.Now().Format("2006-01-02"),
			Created:       time.Now(),
		}
		trades = append(trades, trade)
	} else {
		flag.Usage()
		log.Fatal("Please provide order details via flags or a file path")
	}

	// 4. Upsert into MongoDB
	if len(trades) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Ensure unique index on tracking
		if err := repo.CreateUniqueTrackingIndex(ctx); err != nil {
			log.Printf("Warning: Could not create unique index on tracking: %v", err)
		}

		res, err := repo.BulkUpsertByTracking(ctx, trades)
		if err != nil {
			log.Fatalf("Failed to bulk upsert trades: %v", err)
		}
		fmt.Printf("Successfully processed %d orders (Matched: %d, Upserted: %d, Modified: %d)\n", 
			len(trades), res.MatchedCount, res.UpsertedCount, res.ModifiedCount)
	}
}
