package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// CLI Flags
	collection := flag.String("collection", "sales", "Collection to query (sales or orders)")
	queryJSON := flag.String("query", "{}", `MongoDB query in JSON format.
Examples:
  '{"qty": 10}'
  '{"order_date": "2026-02-01"}'`)
	limit := flag.Int64("limit", 10, "Limit the number of results")
	verbose := flag.Bool("v", false, "Brief output (date, price, tracking)")
	veryVerbose := flag.Bool("vv", false, "Full output (entire record)")
	fields := flag.String("f", "", "Comma-separated list of specific fields to display")
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

	// 3. Parse JSON Query
	var filter map[string]interface{}
	err = json.Unmarshal([]byte(*queryJSON), &filter)
	if err != nil {
		log.Fatalf("Invalid JSON query: %v", err)
	}

	// Helper to convert date strings to time.Time in the filter
	processFilter(filter)

	// 4. Execute Query
	coll := client.Database("tradesmart").Collection(*collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("Error decoding results: %v", err)
	}

	fmt.Printf("Query: %s\n", *queryJSON)
	fmt.Printf("Found %d results.\n\n", len(results))

	// 5. Display Logic
	if *verbose || *veryVerbose || *fields != "" {
		fieldList := []string{}
		if *fields != "" {
			for _, f := range strings.Split(*fields, ",") {
				fieldList = append(fieldList, strings.TrimSpace(f))
			}
		}

		for i, res := range results {
			if int64(i) >= *limit {
				break
			}
			fmt.Printf("[%d] ", i+1)

			if *veryVerbose {
				// Print everything
				for k, v := range res {
					fmt.Printf("%s: %v | ", k, formatValue(v))
				}
			} else if len(fieldList) > 0 {
				// Print specific fields
				for _, f := range fieldList {
					fmt.Printf("%s: %v | ", f, formatValue(res[f]))
				}
			} else if *verbose {
				// Brief output: Date, Price, Tracking
				fmt.Printf("Date: %v | Price: %v | Tracking: %v",
					formatValue(res["order_date"]),
					formatValue(res["sale_price_1"]),
					formatValue(res["tracking"]))
			}
			fmt.Println()
		}
	}
}

func processFilter(m map[string]interface{}) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			if t, err := time.Parse("2006-01-02", val); err == nil {
				m[k] = t
			} else if t, err := time.Parse(time.RFC3339, val); err == nil {
				m[k] = t
			}
		case map[string]interface{}:
			processFilter(val)
		}
	}
}

func formatValue(v interface{}) interface{} {
	if v == nil {
		return "<nil>"
	}
	if t, ok := v.(primitive.DateTime); ok {
		return t.Time().Format("2006-01-02")
	}
	return v
}
