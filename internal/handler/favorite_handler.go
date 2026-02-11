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
)

// AddToFavorites handles POST /api/favorites
func AddToFavorites(favRepo repository.FavoriteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var input model.AddFavoriteInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate car ID
		carID, err := primitive.ObjectIDFromHex(input.CarID)
		if err != nil {
			http.Error(w, "Invalid car ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Add to favorites
		if err := favRepo.AddToFavorites(ctx, userID, carID); err != nil {
			http.Error(w, "Failed to add to favorites", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Added to favorites successfully",
		})
	}
}

// GetFavorites handles GET /api/favorites
func GetFavorites(favRepo repository.FavoriteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		favorites, err := favRepo.GetUserFavorites(ctx, userID)
		if err != nil {
			http.Error(w, "Failed to fetch favorites", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(favorites)
	}
}

// RemoveFromFavorites handles DELETE /api/favorites/{carId}
func RemoveFromFavorites(favRepo repository.FavoriteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract car ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/favorites/")
		carID, err := primitive.ObjectIDFromHex(path)
		if err != nil {
			http.Error(w, "Invalid car ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := favRepo.RemoveFromFavorites(ctx, userID, carID); err != nil {
			http.Error(w, "Failed to remove from favorites", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Removed from favorites successfully",
		})
	}
}

// GetFavoritesCount handles GET /api/favorites/count
func GetFavoritesCount(favRepo repository.FavoriteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := favRepo.GetFavoritesCount(ctx, userID)
		if err != nil {
			http.Error(w, "Failed to get favorites count", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{
			"count": count,
		})
	}
}
