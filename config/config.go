package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	GoogleClientIDs []string
	Port            string
	EmailWhitelist  []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// env file not used in deployed environments, don't error
		// return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		os.Getenv("POSTGRES_USER"),
		"****",
		os.Getenv("POSTGRES_SERVER"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	fmt.Println(url)

	return &Config{
		// DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_SERVER"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB")),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		GoogleClientIDs: parseGoogleClientIDs(os.Getenv("GOOGLE_CLIENT_IDS")),
		// todo: add defaults for some of these
		Port: os.Getenv("PORT"),
		EmailWhitelist: []string{
			"radiatus.io",
			// Add more allowed domains or full email addresses here
		},
	}, nil
}

func parseGoogleClientIDs(envValue string) []string {
	if envValue == "" {
		return nil
	}
	return strings.Split(envValue, ",")
}
