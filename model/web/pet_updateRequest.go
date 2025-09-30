package web

type PetUpdateRequest struct {
	Id int `validate:"required"`
	Name string `validate:"required,min=3,max=100"`
	Species string `validate:"required"`
	Price float64 `validate:"required"`
}