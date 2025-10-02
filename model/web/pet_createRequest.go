package web

type PetCreateRequest struct {
    Name    string  `json:"name" validate:"required"`
    Species string  `json:"species" validate:"required"`
    Price   float64 `json:"price" validate:"gte=0"`
}