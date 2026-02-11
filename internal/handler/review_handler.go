package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/teamserik/online-car-store/internal/middleware"
	"github.com/teamserik/online-car-store/internal/model"
	"github.com/teamserik/online-car-store/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateReview handles POST /api/reviews
func CreateReview(reviewRepo repository.ReviewRepository, userRepo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var carIDStr string
		var rating int
		var comment string

		// Проверяем откуда приходит car_id
		carIDQuery := r.URL.Query().Get("car_id")

		if carIDQuery != "" {
			// car_id пришел через query параметр
			carIDStr = carIDQuery

			// Читаем остальные данные из body
			var input struct {
				Rating  int    `json:"rating"`
				Comment string `json:"comment"`
			}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			rating = input.Rating
			comment = input.Comment
		} else {
			// Читаем все из body
			var input model.CreateReviewInput
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			carIDStr = input.CarID
			rating = input.Rating
			comment = input.Comment
		}

		// Validate rating
		if rating < 1 || rating > 5 {
			http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
			return
		}

		// Parse car ID
		carObjectID, err := primitive.ObjectIDFromHex(carIDStr)
		if err != nil {
			http.Error(w, "Invalid car ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get username
		user, err := userRepo.GetUserByID(ctx, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		review := &model.Review{
			CarID:    carObjectID,
			UserID:   userID,
			Username: user.Username,
			Rating:   rating,
			Comment:  comment,
		}

		if err := reviewRepo.CreateReview(ctx, review); err != nil {
			http.Error(w, "Failed to create review", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(review)
	}
}

// GetCarReviews handles GET /api/reviews?car_id=xxx
func GetCarReviews(reviewRepo repository.ReviewRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get car_id from query parameters
		carID := r.URL.Query().Get("car_id")
		if carID == "" {
			http.Error(w, "car_id parameter is required", http.StatusBadRequest)
			return
		}

		carObjectID, err := primitive.ObjectIDFromHex(carID)
		if err != nil {
			http.Error(w, "Invalid car ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		reviewsResponse, err := reviewRepo.GetCarReviews(ctx, carObjectID)
		if err != nil {
			http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reviewsResponse)
	}
}

// UpdateReview handles PUT /api/reviews/{reviewId}
func UpdateReview(reviewRepo repository.ReviewRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract review ID from URL
		path := strings.TrimPrefix(r.URL.Path, "/api/reviews/")
		reviewID, err := primitive.ObjectIDFromHex(path)
		if err != nil {
			http.Error(w, "Invalid review ID", http.StatusBadRequest)
			return
		}

		var input model.UpdateReviewInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate rating
		if input.Rating < 1 || input.Rating > 5 {
			http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = reviewRepo.UpdateReview(ctx, reviewID, userID, input)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Review not found or you don't have permission", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update review", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Review updated successfully",
		})
	}
}

// DeleteReview handles DELETE /api/reviews/{reviewId}
func DeleteReview(reviewRepo repository.ReviewRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract review ID from URL
		path := strings.TrimPrefix(r.URL.Path, "/api/reviews/")
		reviewID, err := primitive.ObjectIDFromHex(path)
		if err != nil {
			http.Error(w, "Invalid review ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = reviewRepo.DeleteReview(ctx, reviewID, userID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Review not found or you don't have permission", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete review", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Review deleted successfully",
		})
	}
}
