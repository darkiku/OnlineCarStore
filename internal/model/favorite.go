package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Favorite represents a user's favorite car
type Favorite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	CarID     primitive.ObjectID `bson:"car_id" json:"car_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// FavoriteWithCar represents a favorite with full car details
type FavoriteWithCar struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	CarID     primitive.ObjectID `json:"car_id"`
	Car       *Car               `json:"car"`
	CreatedAt time.Time          `json:"created_at"`
}

// AddFavoriteInput for adding a car to favorites
type AddFavoriteInput struct {
	CarID string `json:"car_id"`
}
