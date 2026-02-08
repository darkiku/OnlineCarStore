package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Car struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Make         string             `bson:"make" json:"make"`
	Model        string             `bson:"model" json:"model"`
	Year         int                `bson:"year" json:"year"`
	Price        float64            `bson:"price" json:"price"`
	Mileage      int                `bson:"mileage" json:"mileage"`
	BodyType     string             `bson:"body_type" json:"body_type"`
	FuelType     string             `bson:"fuel_type" json:"fuel_type"`       // Бензин, Дизель, Электро
	Transmission string             `bson:"transmission" json:"transmission"` // Автомат, Механика
	Color        string             `bson:"color" json:"color"`
	HorsePower   int                `bson:"horsepower" json:"horsepower"`
	EngineSize   float64            `bson:"engine_size" json:"engine_size"` // литры
	Description  string             `bson:"description" json:"description"`
	ImageURL     string             `bson:"image_url" json:"image_url"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateCarInput struct {
	Make         string  `json:"make"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	Price        float64 `json:"price"`
	Mileage      int     `json:"mileage"`
	BodyType     string  `json:"body_type"`
	FuelType     string  `json:"fuel_type"`
	Transmission string  `json:"transmission"`
	Color        string  `json:"color"`
	HorsePower   int     `json:"horsepower"`
	EngineSize   float64 `json:"engine_size"`
	Description  string  `json:"description"`
	ImageURL     string  `json:"image_url"`
}

type UpdateCarInput struct {
	Price       *float64 `json:"price,omitempty"`
	Mileage     *int     `json:"mileage,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type FilterParams struct {
	MinPrice     *float64 `json:"min_price"`
	MaxPrice     *float64 `json:"max_price"`
	Make         *string  `json:"make"`
	BodyType     *string  `json:"body_type"`
	FuelType     *string  `json:"fuel_type"`
	Transmission *string  `json:"transmission"`
	MinYear      *int     `json:"min_year"`
	MaxYear      *int     `json:"max_year"`
}
