package repository

import (
	"context"
	"errors"
	"time"

	"github.com/teamserik/online-car-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, id primitive.ObjectID, user *model.User) error
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &MongoUserRepository{
		collection: collection,
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Role == "" {
		user.Role = "user"
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, id primitive.ObjectID, user *model.User) error {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"phone":      user.Phone,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserByID is an alias for FindByID (for consistency with review handler)
func (r *MongoUserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return r.FindByID(ctx, id)
}
