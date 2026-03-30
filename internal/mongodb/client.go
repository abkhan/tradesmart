package mongodb

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"mongotest/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps the mongo.Client to provide convenience methods.
type Client struct {
	*mongo.Client
	isHealthy bool
}

// Connect initializes a new MongoDB client using the provided config.
func Connect(cfg *config.Config) (*Client, error) {
	// Standard Atlas SRV URI construction
	uri := fmt.Sprintf("%s://%s:%s@%s%s",
		cfg.MongoScheme, cfg.MongoUser, cfg.MongoPassword, cfg.MongoHost, cfg.MongoURI)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false, // Standard practice; Atlas handles TLS
		})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	c := &Client{Client: client, isHealthy: true}
	return c, nil
}

// StartHealthCheck starts a background goroutine to monitor MongoDB connectivity.
func (c *Client) StartHealthCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := c.Ping(ctx)
			cancel()

			c.isHealthy = (err == nil)
			if err != nil {
				log.Printf("Health check failed: %v", err)
			}

			<-ticker.C
		}
	}()
}

// IsHealthy returns the current connectivity state.
func (c *Client) IsHealthy() bool {
	return c.isHealthy
}

// Disconnect closes the MongoDB connection.
func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.Client.Disconnect(ctx)
}

// Ping sends a ping to the database.
func (c *Client) Ping(ctx context.Context) error {
	return c.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err()
}
