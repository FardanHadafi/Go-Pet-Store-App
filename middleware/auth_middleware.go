package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"Go-PetStoreApp/helper"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

type JWTMiddleware struct{}

func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{}
}

func (m *JWTMiddleware) Authenticate(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		claims, err := helper.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next(w, r.WithContext(ctx), ps)
	}
}

// RequireRole returns a wrapper that requires a specific role (e.g., "admin")
func (m *JWTMiddleware) RequireRole(role string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		roleVal, ok := r.Context().Value(RoleKey).(string)
		if !ok || roleVal != role {
			writeError(w, "forbidden: insufficient role", http.StatusForbidden)
			return
		}
		next(w, r, ps)
	}
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(UserIDKey).(int)
	return v, ok
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(EmailKey).(string)
	return v, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(RoleKey).(string)
	return v, ok
}
