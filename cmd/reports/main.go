package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// CLI Flags
	reportType := flag.String("type", "summary", "Type of report: summary, daily, pending, returns")
	collection := flag.String("collection", "orders", "Collection to run report on (orders is default)")
	dateStr := flag.String("date", time.Now().Format("2006-01-02"), "Date for daily report (YYYY-MM-DD)")
	monthStr := flag.String("month", time.Now().Format("2006-01"), "Month for reports (YYYY-MM)")
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

	coll := client.Database("tradesmart").Collection(*collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. Routing Report Logic
	switch *reportType {
	case "daily":
		runDailyReport(ctx, coll, *dateStr)
	case "pending":
		runPendingReport(ctx, coll, *monthStr)
	case "returns":
		runReturnsReport(ctx, coll, *monthStr)
	case "summary":
		runFullSummary(ctx, coll, *monthStr)
	default:
		fmt.Printf("Unknown report type: %s\n", *reportType)
	}
}

func runDailyReport(ctx context.Context, coll *mongo.Collection, dateStr string) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Invalid date format %s: %v", dateStr, err)
	}
	filter := bson.M{
		"order_date": t,
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	var results []bson.M
	cursor.All(ctx, &results)

	var totalSales, totalProfit float64
	pendingCount := 0

	for _, res := range results {
		totalSales += toFloat(res["sale_price_1"])
		totalProfit += toFloat(res["profit"])
		if toString(res["tracking"]) == "" {
			pendingCount++
		}
	}

	fmt.Printf("--- Daily Report (%s): %s ---\n", coll.Name(), dateStr)
	fmt.Printf("Total Orders:  %d\n", len(results))
	fmt.Printf("Total Sales:   $%.2f\n", totalSales)
	fmt.Printf("Total Profit:  $%.2f\n", totalProfit)
	fmt.Printf("Pending (No Tracking): %d\n", pendingCount)
}

func runPendingReport(ctx context.Context, coll *mongo.Collection, monthStr string) {
	start, err := time.Parse("2006-01", monthStr)
	if err != nil {
		log.Fatalf("Invalid month format %s: %v", monthStr, err)
	}
	end := start.AddDate(0, 1, 0)

	filter := bson.M{
		"order_date": bson.M{"$gte": start, "$lt": end},
		"tracking":   "",
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	var results []bson.M
	cursor.All(ctx, &results)

	fmt.Printf("--- Pending Orders (No Tracking) (%s) for %s ---\n", coll.Name(), monthStr)
	fmt.Printf("Count: %d\n", len(results))
	for i, res := range results {
		fmt.Printf("[%d] ID: %v | Type: %v | Refund: %v\n", 
			i+1, res["seller_order_id"], res["order_type"], res["return_refund"])
	}
}

func runReturnsReport(ctx context.Context, coll *mongo.Collection, monthStr string) {
	start, err := time.Parse("2006-01", monthStr)
	if err != nil {
		log.Fatalf("Invalid month format %s: %v", monthStr, err)
	}
	end := start.AddDate(0, 1, 0)

	filter := bson.M{
		"order_date":    bson.M{"$gte": start, "$lt": end},
		"return_refund": bson.M{"$gt": 0},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	var results []bson.M
	cursor.All(ctx, &results)

	var totalReturned float64
	for _, res := range results {
		totalReturned += toFloat(res["return_refund"])
	}

	fmt.Printf("--- Returns Report (%s) for %s ---\n", coll.Name(), monthStr)
	fmt.Printf("Total Returned Orders: %d\n", len(results))
	fmt.Printf("Total Return Value:    $%.2f\n", totalReturned)
}

func runFullSummary(ctx context.Context, coll *mongo.Collection, monthStr string) {
	start, err := time.Parse("2006-01", monthStr)
	if err != nil {
		log.Fatalf("Invalid month format %s: %v", monthStr, err)
	}
	end := start.AddDate(0, 1, 0)

	filter := bson.M{
		"order_date": bson.M{"$gte": start, "$lt": end},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	var results []bson.M
	cursor.All(ctx, &results)

	var sales, profit, returns, net float64
	for _, res := range results {
		sales += toFloat(res["sale_price_1"])
		profit += toFloat(res["profit"])
		returns += toFloat(res["return_refund"])
		net += toFloat(res["net"])
	}

	fmt.Printf("--- Full Summary (%s) for %s ---\n", coll.Name(), monthStr)
	fmt.Printf("Total Records:  %d\n", len(results))
	fmt.Printf("Gross Sales:    $%.2f\n", sales)
	fmt.Printf("Total Returns:  $%.2f\n", returns)
	fmt.Printf("Total Profit:   $%.2f\n", profit)
	fmt.Printf("Net Total:      $%.2f\n", net)
}

func toFloat(v interface{}) float64 {
	if v == nil { return 0 }
	switch val := v.(type) {
	case float64: return val
	case int32: return float64(val)
	case int64: return float64(val)
	}
	return 0
}

func toString(v interface{}) string {
	if v == nil { return "" }
	if s, ok := v.(string); ok { return s }
	return ""
}
