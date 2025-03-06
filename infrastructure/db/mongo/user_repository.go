package mongo

import (
	"context"
	"errors"

	"github.com/pedro00627/urblog/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Save(user *domain.User) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"id": user.ID},
		bson.M{"$set": user},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *UserRepository) FindByID(userID string) (*domain.User, error) {
	filter := bson.M{"id": userID}
	var user domain.User
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepository) FindByName(s string) (*domain.User, error) {
	filter := bson.M{"username": s}
	var user domain.User
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}
