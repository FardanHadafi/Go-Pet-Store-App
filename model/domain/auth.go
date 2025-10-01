package domain

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email,omitempty"`
}

type AuthResponse struct {
    Token   string `json:"token"`
    Username string `json:"username"`
    Message string `json:"message"`
}