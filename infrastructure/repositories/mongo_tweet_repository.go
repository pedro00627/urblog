package repositories

import (
	"context"

	"github.com/pedro00627/urblog/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTweetRepository struct {
	collection *mongo.Collection
}

func NewMongoTweetRepository(db *mongo.Database) *MongoTweetRepository {
	return &MongoTweetRepository{
		collection: db.Collection("tweets"),
	}
}

func (r *MongoTweetRepository) Save(tweet *domain.Tweet) error {
	_, err := r.collection.InsertOne(context.TODO(), tweet)
	return err
}

func (r *MongoTweetRepository) FindByUserID(userID string, limit, offset int) ([]*domain.Tweet, error) {
	filter := bson.M{"userid": userID}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tweets []*domain.Tweet
	for cursor.Next(context.TODO()) {
		var tweet domain.Tweet
		err := cursor.Decode(&tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}
