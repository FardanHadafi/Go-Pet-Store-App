package web

type PetResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Species string `json:"species"`
	Price float64 `json:"price"`
}