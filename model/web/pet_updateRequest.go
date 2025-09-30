package web

type PetUpdateRequest struct {
	Id int
	Name string
	Species string
	Price float64
}