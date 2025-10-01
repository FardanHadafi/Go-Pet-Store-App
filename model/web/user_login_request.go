package web

type UserLoginRequest struct {
    Username string `validate:"required,min=3"`
    Password string `validate:"required" json:"password"`
}