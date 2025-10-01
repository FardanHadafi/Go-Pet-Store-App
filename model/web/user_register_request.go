package web

type UserRegisterRequest struct {
    Username string `validate:"required,min=3,max=50" json:"username"`
    Email    string `validate:"required,email" json:"email"`
    Password string `validate:"required,min=8" json:"password"`
}