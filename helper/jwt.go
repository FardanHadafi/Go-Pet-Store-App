package helper

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID int, email string) (string, error) {
    // Get secret from environment variable
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "My-Secret-Key" // Default for development
    }

    // Create claims
    claims := JWTClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign token with secret
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "My-Secret-Key"
    }

    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    // Extract claims
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// GetUserIDFromToken extracts user ID from token string
func GetUserIDFromToken(tokenString string) (int, error) {
    claims, err := ValidateToken(tokenString)
    if err != nil {
        return 0, err
    }
    return claims.UserID, nil
}