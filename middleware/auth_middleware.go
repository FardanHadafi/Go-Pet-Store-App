package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

// JWTClaims defines the structure of JWT claims
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// EmailKey is the context key for email
	EmailKey contextKey = "email"
)

// JWTMiddleware creates a middleware for JWT authentication
type JWTMiddleware struct {
	SecretKey string
}

// NewJWTMiddleware creates a new JWT middleware instance
// auth_middleware.go
func NewJWTMiddleware() *JWTMiddleware {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    return &JWTMiddleware{ SecretKey: secretKey }
}


// Authenticate is the middleware handler that validates JWT tokens
func (m *JWTMiddleware) Authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.writeErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.writeErrorResponse(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.SecretKey), nil
		})

		if err != nil {
			m.writeErrorResponse(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			m.writeErrorResponse(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		// Call the next handler with updated context
		next(w, r.WithContext(ctx), ps)
	}
}

// OptionalAuthenticate is middleware that extracts JWT if present but doesn't require it
func (m *JWTMiddleware) OptionalAuthenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				
				token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(m.SecretKey), nil
				})

				if err == nil {
					if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
						ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
						ctx = context.WithValue(ctx, EmailKey, claims.Email)
						r = r.WithContext(ctx)
					}
				}
			}
		}

		next(w, r, ps)
	}
}

// writeErrorResponse writes a JSON error response
func (m *JWTMiddleware) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}