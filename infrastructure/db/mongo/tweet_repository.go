package mongo

import (
	"context"

	"github.com/pedro00627/urblog/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TweetRepository struct {
	collection *mongo.Collection
}

func NewTweetRepository(db *mongo.Database) *TweetRepository {
	return &TweetRepository{
		collection: db.Collection("tweets"),
	}
}

func (r *TweetRepository) Save(tweet *domain.Tweet) error {
	_, err := r.collection.InsertOne(context.TODO(), tweet)
	return err
}

func (r *TweetRepository) FindByUserID(userID string, limit, offset int) ([]*domain.Tweet, error) {
	filter := bson.M{"userid": userID}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		closeErr := cursor.Close(ctx)
		if closeErr != nil {

		}
	}(cursor, context.TODO())

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
