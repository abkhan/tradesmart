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
	r.HandleFunc("/api/orders/search", searchOrders).Methods("GET")
	r.HandleFunc("/api/orders/{orderId}", getOrderByID).Methods("GET")
	r.HandleFunc("/api/orders/tracking/{trackingId}", getOrderByTracking).Methods("GET")

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

	log.Println("Orders API starting on :8080 with CORS allowed for 10.0.0.206 and 10.0.0.164")
	log.Fatal(srv.ListenAndServe())
}

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trade, err := repo.GetByOrderID(ctx, orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trade)
}

func getOrderByTracking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackingID := vars["trackingId"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trade, err := repo.GetByTracking(ctx, trackingID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trade)
}

func searchOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trades, err := repo.Search(ctx, query)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}
