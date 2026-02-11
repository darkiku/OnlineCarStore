package repository

import (
	"context"
	"time"

	"github.com/teamserik/online-car-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FavoriteRepository interface {
	AddToFavorites(ctx context.Context, userID, carID primitive.ObjectID) error
	RemoveFromFavorites(ctx context.Context, userID, carID primitive.ObjectID) error
	GetUserFavorites(ctx context.Context, userID primitive.ObjectID) ([]model.FavoriteWithCar, error)
	GetFavoritesCount(ctx context.Context, userID primitive.ObjectID) (int64, error)
	IsFavorite(ctx context.Context, userID, carID primitive.ObjectID) (bool, error)
}

type MongoFavoriteRepository struct {
	collection     *mongo.Collection
	carsCollection *mongo.Collection
}

func NewMongoFavoriteRepository(collection *mongo.Collection, carsCollection *mongo.Collection) *MongoFavoriteRepository {
	return &MongoFavoriteRepository{
		collection:     collection,
		carsCollection: carsCollection,
	}
}

// AddToFavorites adds a car to user's favorites
func (r *MongoFavoriteRepository) AddToFavorites(ctx context.Context, userID, carID primitive.ObjectID) error {
	// Check if already in favorites
	exists, err := r.IsFavorite(ctx, userID, carID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already in favorites
	}

	favorite := model.Favorite{
		UserID:    userID,
		CarID:     carID,
		CreatedAt: time.Now(),
	}

	_, err = r.collection.InsertOne(ctx, favorite)
	return err
}

// RemoveFromFavorites removes a car from user's favorites
func (r *MongoFavoriteRepository) RemoveFromFavorites(ctx context.Context, userID, carID primitive.ObjectID) error {
	filter := bson.M{
		"user_id": userID,
		"car_id":  carID,
	}

	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// GetUserFavorites returns all favorites for a user with car details
func (r *MongoFavoriteRepository) GetUserFavorites(ctx context.Context, userID primitive.ObjectID) ([]model.FavoriteWithCar, error) {
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var favorites []model.Favorite
	if err = cursor.All(ctx, &favorites); err != nil {
		return nil, err
	}

	// Fetch car details for each favorite
	var favoritesWithCars []model.FavoriteWithCar
	for _, fav := range favorites {
		var car model.Car
		err := r.carsCollection.FindOne(ctx, bson.M{"_id": fav.CarID}).Decode(&car)
		if err != nil {
			// If car not found, skip it
			if err == mongo.ErrNoDocuments {
				continue
			}
			return nil, err
		}

		favoritesWithCars = append(favoritesWithCars, model.FavoriteWithCar{
			ID:        fav.ID,
			UserID:    fav.UserID,
			CarID:     fav.CarID,
			Car:       &car,
			CreatedAt: fav.CreatedAt,
		})
	}

	return favoritesWithCars, nil
}

// GetFavoritesCount returns the count of favorites for a user
func (r *MongoFavoriteRepository) GetFavoritesCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	filter := bson.M{"user_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}

// IsFavorite checks if a car is in user's favorites
func (r *MongoFavoriteRepository) IsFavorite(ctx context.Context, userID, carID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"user_id": userID,
		"car_id":  carID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
