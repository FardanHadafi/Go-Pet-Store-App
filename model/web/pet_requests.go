package web

type PetCreateRequest struct {
	Name    string  `json:"name" validate:"required"`
	Species string  `json:"species" validate:"required"`
	Price   float64 `json:"price" validate:"required,gte=0"`
}

type PetUpdateRequest struct {
	Id      int     `json:"id"`
	Name    string  `json:"name" validate:"required"`
	Species string  `json:"species" validate:"required"`
	Price   float64 `json:"price" validate:"required,gte=0"`
}
