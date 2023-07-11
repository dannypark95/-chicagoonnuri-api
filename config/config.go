package config

import (
	"os"
	"sync"
)

var (
	once sync.Once

	// DatabaseURL is the MongoDB connection URL.
	DatabaseURL string

	// AWSAccessKey is the AWS access key for your account.
	AWSAccessKey string

	// AWSSecretKey is the AWS secret key for your account.
	AWSSecretKey string

	// JWTSecret is the secret key used for JWT authentication.
	JWTSecret string
)

// LoadEnv loads environment variables into the config package variables.
func LoadEnv() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	AWSAccessKey = os.Getenv("AWS_ACCESS_KEY")
	AWSSecretKey = os.Getenv("AWS_SECRET_KEY")
	JWTSecret = os.Getenv("JWT_SECRET")
}

// func init() {
// 	once.Do(LoadEnv)
// 	if AWSAccessKey == "" || AWSSecretKey == "" {
// 		log.Fatal("AWS credentials are not set in the environment.")
// 	}
// }
