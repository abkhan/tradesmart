package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"mongotest/internal/parser"
)

func main() {
	file1 := flag.String("f1", "", "First file to compare")
	file2 := flag.String("f2", "", "Second file to compare")
	verbose := flag.Bool("v", false, "Show overlapping Order IDs")
	flag.Parse()

	if *file1 == "" || *file2 == "" {
		log.Fatal("Please provide two files to compare using -f1 and -f2")
	}

	// 1. Load data from both files (ignoring dates for diff)
	trades1, err := parser.LoadTradesFromFile(*file1, "")
	if err != nil {
		log.Fatalf("Error loading file 1: %v", err)
	}

	trades2, err := parser.LoadTradesFromFile(*file2, "")
	if err != nil {
		log.Fatalf("Error loading file 2: %v", err)
	}

	// 2. Index first file by Tracking and OrderID
	trackingMap := make(map[string]bool)
	orderIDMap := make(map[string]bool)

	for _, t := range trades1 {
		if t.Tracking != "" {
			trackingMap[strings.TrimSpace(t.Tracking)] = true
		}
		if t.OrderID != "" {
			orderIDMap[strings.TrimSpace(t.OrderID)] = true
		}
	}

	// 3. Check for duplicates in the second file
	duplicateTracking := make(map[string]bool)
	duplicateOrderIDs := make(map[string]bool)

	for _, t := range trades2 {
		track := strings.TrimSpace(t.Tracking)
		if track != "" && trackingMap[track] {
			duplicateTracking[track] = true
		}

		oid := strings.TrimSpace(t.OrderID)
		if oid != "" && orderIDMap[oid] {
			duplicateOrderIDs[oid] = true
		}
	}

	// 4. Output results
	fmt.Printf("File 1: %s (%d records)\n", *file1, len(trades1))
	fmt.Printf("File 2: %s (%d records)\n", *file2, len(trades2))
	fmt.Printf("Overlap Summary:\n")
	fmt.Printf("  - Overlapping Tracking Numbers: %d\n", len(duplicateTracking))
	fmt.Printf("  - Overlapping Order IDs: %d\n", len(duplicateOrderIDs))

	if *verbose {
		if len(duplicateOrderIDs) > 0 {
			fmt.Println("\nOverlapping Order IDs:")
			for oid := range duplicateOrderIDs {
				fmt.Printf("  %s\n", oid)
			}
		}
		if len(duplicateTracking) > 0 {
			fmt.Println("\nOverlapping Tracking Numbers:")
			for track := range duplicateTracking {
				fmt.Printf("  %s\n", track)
			}
		}
	}
}
