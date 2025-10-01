package domain

type Pet struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Species string  `json:"species"`
	Price   float64 `json:"price"`
}