package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"mongotest/internal/config"
	"mongotest/internal/mongodb"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var repo *mongodb.Repository

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect()

	repo = mongodb.NewRepository(client.Database("tradesmart"), "orders")

	r := mux.NewRouter()
	r.HandleFunc("/api/reports/summary", getSummary).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://10.0.0.206", "http://10.0.0.164", "http://localhost:3000", "http://10.0.0.206:*", "http://10.0.0.164:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Reports API starting on :8080 with CORS allowed for 10.0.0.206 and 10.0.0.164")
	log.Fatal(srv.ListenAndServe())
}

func getSummary(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("startDate")
	endStr := r.URL.Query().Get("endDate")

	if startStr == "" || endStr == "" {
		http.Error(w, "Query parameters 'startDate' and 'endDate' (YYYY-MM-DD) are required", http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		http.Error(w, "Invalid startDate format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		http.Error(w, "Invalid endDate format", http.StatusBadRequest)
		return
	}

	// Make end date inclusive of that day
	end = end.AddDate(0, 0, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	summary, err := repo.GetSummary(ctx, start, end)
	if err != nil {
		http.Error(w, "Failed to generate summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
