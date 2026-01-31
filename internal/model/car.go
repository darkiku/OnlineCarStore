package model

import "time"

type Car struct {
	ID          int       `json:"id"`
	Make        string    `json:"make"` // Toyota, BMW и т.д.
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	Price       float64   `json:"price"`
	Mileage     int       `json:"mileage,omitempty"`   // пробег, км
	BodyType    string    `json:"body_type,omitempty"` // sedan, suv, hatchback...
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateCarInput struct {
	Make        string  `json:"make"`
	Model       string  `json:"model"`
	Year        int     `json:"year"`
	Price       float64 `json:"price"`
	Mileage     int     `json:"mileage,omitempty"`
	BodyType    string  `json:"body_type,omitempty"`
	Description string  `json:"description,omitempty"`
}

type UpdateCarInput struct {
	Price       *float64 `json:"price,omitempty"`
	Mileage     *int     `json:"mileage,omitempty"`
	Description *string  `json:"description,omitempty"`
}
