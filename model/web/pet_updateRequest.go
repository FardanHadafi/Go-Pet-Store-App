package web

type PetUpdateRequest struct {
	Id int `validate:"required" json:"id"`
	Name string `validate:"required,min=3,max=100" json:"name"`
	Species string `validate:"required" json:"species"`
	Price float64 `validate:"required" json:"price"`
}