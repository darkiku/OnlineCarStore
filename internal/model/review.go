package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a car review
type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CarID     primitive.ObjectID `bson:"car_id" json:"car_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Rating    int                `bson:"rating" json:"rating"` // 1-5 stars
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateReviewInput for adding a review
type CreateReviewInput struct {
	CarID   string `json:"car_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// UpdateReviewInput for updating a review
type UpdateReviewInput struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// ReviewsResponse with aggregated data
type ReviewsResponse struct {
	Reviews       []Review `json:"reviews"`
	AverageRating float64  `json:"average_rating"`
	TotalReviews  int      `json:"total_reviews"`
}
