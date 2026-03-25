package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MongoUser     string `json:"mongo_user"`
	MongoPassword string `json:"mongo_password"`
	MongoScheme   string `json:"mongo_scheme"`
	MongoHost     string `json:"mongo_host"`
	MongoURI      string `json:"mongo_uri"`
}

func LoadConfig(path string) (*Config, error) {
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
