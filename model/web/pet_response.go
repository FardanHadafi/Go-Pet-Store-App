package web

import "time"

type PetResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Species   string    `json:"species"`
	Price     float64   `json:"price"`
	OwnerId   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
