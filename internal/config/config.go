package config

import (
	"fmt"
	"os"
)

// Config defines the application's configuration.
type Config struct {
	DB DBConfig
	// other configuration
}

// DBConfig holds database configuration.
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// LoadConfig loads the application's configuration
func LoadConfig() (Config, error) {
	cfg := Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		// Load other configuration values
	}
	//check if required configs are set
	if cfg.DB.Host == "" || cfg.DB.Port == "" || cfg.DB.Name == "" || cfg.DB.User == "" || cfg.DB.Password == "" {
		return Config{}, fmt.Errorf("database configuration is incomplete")
	}
	return cfg, nil
}
