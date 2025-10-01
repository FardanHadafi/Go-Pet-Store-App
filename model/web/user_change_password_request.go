package web

type UserChangePasswordRequest struct {
    Id          int    `validate:"required"`
    OldPassword string `validate:"required" json:"old_password"`
    NewPassword string `validate:"required,min=8" json:"new_password"`
}