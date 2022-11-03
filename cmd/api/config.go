package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	MongoDBConnectionString string `envconfig:"mongo_connection_string" required:"true"`
	CacheSize               int    `envconfig:"cache_size" required:"false" default:"5000"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
