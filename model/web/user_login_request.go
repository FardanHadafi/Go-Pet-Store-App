package web

type UserLoginRequest struct {
    Username string `validate:"required,min=3"`
    Email    string `validate:"required,email" json:"email"`
    Password string `validate:"required" json:"password"`
}