package app

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	JWTSecretKey string
	TokenExpiry  int // in hours
}

// LoadConfig loads environment variables into Config struct
func LoadConfig() *Config {
	_ = godotenv.Load() // don’t panic if .env not found, system vars might exist

	// Token Expiry
	expiryStr := os.Getenv("TOKEN_EXPIRATION_HOURS")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil || expiry <= 0 {
		log.Println("⚠️ TOKEN_EXPIRATION_HOURS not set or invalid, defaulting to 24h")
		expiry = 24
	}

	// JWT Secret
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("❌ JWT_SECRET_KEY must be set in environment variables")
	}

	return &Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBSSLMode:    os.Getenv("DB_SSLMODE"),
		JWTSecretKey: secret,
		TokenExpiry:  expiry,
	}
}
