package repository

import (
	"context"
	"time"

	"github.com/teamserik/online-car-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *model.Review) error
	GetCarReviews(ctx context.Context, carID primitive.ObjectID) (*model.ReviewsResponse, error)
	UpdateReview(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID, input model.UpdateReviewInput) error
	DeleteReview(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID) error
	GetReviewByID(ctx context.Context, reviewID primitive.ObjectID) (*model.Review, error)
}

type MongoReviewRepository struct {
	collection *mongo.Collection
}

func NewMongoReviewRepository(collection *mongo.Collection) *MongoReviewRepository {
	return &MongoReviewRepository{
		collection: collection,
	}
}

// CreateReview creates a new review
func (r *MongoReviewRepository) CreateReview(ctx context.Context, review *model.Review) error {
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}

	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetCarReviews returns all reviews for a car with aggregated data
func (r *MongoReviewRepository) GetCarReviews(ctx context.Context, carID primitive.ObjectID) (*model.ReviewsResponse, error) {
	filter := bson.M{"car_id": carID}

	// Sort by created_at descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []model.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	// Calculate average rating
	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}

	averageRating := 0.0
	if len(reviews) > 0 {
		averageRating = float64(totalRating) / float64(len(reviews))
	}

	return &model.ReviewsResponse{
		Reviews:       reviews,
		AverageRating: averageRating,
		TotalReviews:  len(reviews),
	}, nil
}

// UpdateReview updates an existing review (only by the owner)
func (r *MongoReviewRepository) UpdateReview(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID, input model.UpdateReviewInput) error {
	filter := bson.M{
		"_id":     reviewID,
		"user_id": userID, // Ensure user owns the review
	}

	update := bson.M{
		"$set": bson.M{
			"rating":     input.Rating,
			"comment":    input.Comment,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteReview deletes a review (only by the owner)
func (r *MongoReviewRepository) DeleteReview(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":     reviewID,
		"user_id": userID, // Ensure user owns the review
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetReviewByID returns a review by ID
func (r *MongoReviewRepository) GetReviewByID(ctx context.Context, reviewID primitive.ObjectID) (*model.Review, error) {
	var review model.Review
	err := r.collection.FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}
