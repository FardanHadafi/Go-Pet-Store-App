package domain

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email,omitempty"`
}