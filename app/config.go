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
	TokenExpiry  int // hours
}

// LoadConfig loads .env and env vars; fails fast if JWT secret missing.
func LoadConfig() *Config {
	_ = godotenv.Load()

	expiry := 24
	if s := os.Getenv("TOKEN_EXPIRATION_HOURS"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			expiry = v
		} else {
			log.Println("invalid TOKEN_EXPIRATION_HOURS, defaulting to 24")
		}
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY must be set")
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
