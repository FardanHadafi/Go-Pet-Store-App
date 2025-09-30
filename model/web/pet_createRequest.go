package web

type PetCreateRequest struct {
	Name string `validate:"required,min=3,max=100"`
	Species string `validate:"required"`
	Price float64 `validate:"required"`
}