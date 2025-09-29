// model/domain/user.go
package domain

type Pet struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Species string  `json:"species"`
	Price   float64 `json:"price"`
}


type User struct {
    ID           int    `json:"id"`
    Username     string `json:"username"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"`
}

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