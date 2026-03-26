package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	MongoUser      string   `json:"mongo_user"`
	MongoPassword  string   `json:"mongo_password"`
	MongoScheme    string   `json:"mongo_scheme"`
	MongoHost      string   `json:"mongo_host"`
	MongoURI       string   `json:"mongo_uri"`
	AllowedOrigins []string `json:"allowed_origins"`
}

func LoadConfig(path string) (*Config, error) {
	// 1. Try to load from environment variables first (Good for Docker/K8s)
	if os.Getenv("MONGO_USER") != "" {
		origins := []string{"*"}
		if os.Getenv("ALLOWED_ORIGINS") != "" {
			origins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
		}
		return &Config{
			MongoUser:      os.Getenv("MONGO_USER"),
			MongoPassword:  os.Getenv("MONGO_PASSWORD"),
			MongoScheme:    getEnv("MONGO_SCHEME", "mongodb+srv"),
			MongoHost:      getEnv("MONGO_HOST", "cluster0.dbpelmw.mongodb.net"),
			MongoURI:       getEnv("MONGO_URI", "/?appName=Cluster0"),
			AllowedOrigins: origins,
		}, nil
	}

	// 2. Fallback to JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}
	return &cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
