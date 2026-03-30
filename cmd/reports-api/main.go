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
var dbClient *mongodb.Client

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbClient, err = mongodb.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer dbClient.Disconnect()

	// Start background health check (every 1 minute)
	dbClient.StartHealthCheck(1 * time.Minute)

	repo = mongodb.NewRepository(dbClient.Database("tradesmart"), "orders")

	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/reports/summary", getSummary).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
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

	log.Println("Reports API starting on :8080 with health checks enabled")
	log.Fatal(srv.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if dbClient.IsHealthy() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("MongoDB connection lost"))
	}
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
