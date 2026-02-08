package repository

import (
	"context"
	"time"

	"github.com/teamserik/online-car-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	GetByID(ctx context.Context, id string) (*model.Car, error)
	List(ctx context.Context, filter *model.FilterParams) ([]*model.Car, error)
	Update(ctx context.Context, id string, input model.UpdateCarInput) error
	Delete(ctx context.Context, id string) error
}

type mongoCarRepository struct {
	collection *mongo.Collection
}

func NewMongoCarRepository(collection *mongo.Collection) CarRepository {
	return &mongoCarRepository{
		collection: collection,
	}
}

func (r *mongoCarRepository) Create(ctx context.Context, car *model.Car) error {
	car.ID = primitive.NewObjectID()
	car.CreatedAt = time.Now()
	car.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, car)
	return err
}

func (r *mongoCarRepository) GetByID(ctx context.Context, id string) (*model.Car, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var car model.Car
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

func (r *mongoCarRepository) List(ctx context.Context, filter *model.FilterParams) ([]*model.Car, error) {
	query := bson.M{}

	if filter != nil {
		if filter.MinPrice != nil || filter.MaxPrice != nil {
			priceFilter := bson.M{}
			if filter.MinPrice != nil {
				priceFilter["$gte"] = *filter.MinPrice
			}
			if filter.MaxPrice != nil {
				priceFilter["$lte"] = *filter.MaxPrice
			}
			query["price"] = priceFilter
		}

		if filter.Make != nil {
			query["make"] = *filter.Make
		}

		if filter.BodyType != nil {
			query["body_type"] = *filter.BodyType
		}

		if filter.FuelType != nil {
			query["fuel_type"] = *filter.FuelType
		}

		if filter.Transmission != nil {
			query["transmission"] = *filter.Transmission
		}

		if filter.MinYear != nil || filter.MaxYear != nil {
			yearFilter := bson.M{}
			if filter.MinYear != nil {
				yearFilter["$gte"] = *filter.MinYear
			}
			if filter.MaxYear != nil {
				yearFilter["$lte"] = *filter.MaxYear
			}
			query["year"] = yearFilter
		}
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []*model.Car
	if err = cursor.All(ctx, &cars); err != nil {
		return nil, err
	}

	return cars, nil
}

func (r *mongoCarRepository) Update(ctx context.Context, id string, input model.UpdateCarInput) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if input.Price != nil {
		update["$set"].(bson.M)["price"] = *input.Price
	}
	if input.Mileage != nil {
		update["$set"].(bson.M)["mileage"] = *input.Mileage
	}
	if input.Description != nil {
		update["$set"].(bson.M)["description"] = *input.Description
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *mongoCarRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
