package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DatabaseURL string
	DatabaseName string
	JWTSecret   []byte
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DatabaseURL = os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required in .env")
	}

	DatabaseName = os.Getenv("DATABASE_NAME")
	if DatabaseName == "" {
		log.Fatal("DATABASE_NAME is required in .env")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required in .env")
	}
	JWTSecret = []byte(secret)
}
