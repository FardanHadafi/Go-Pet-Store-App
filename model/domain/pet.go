package domain

import "time"

type Pet struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Species   string  `json:"species"`
	Price     float64 `json:"price"`
	CreatedBy int     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}