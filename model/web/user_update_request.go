package web

type UserUpdateRequest struct {
    Id       int    `validate:"required"`
    Username string `validate:"required,min=3,max=50" json:"username"`
    Email    string `validate:"required,email" json:"email"`
}