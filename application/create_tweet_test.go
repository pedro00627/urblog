package application

import (
	"testing"
	"time"

	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
)

type MockTweetRepository struct {
	tweets map[string]*domain.Tweet
}

func NewMockTweetRepository() *MockTweetRepository {
	return &MockTweetRepository{
		tweets: make(map[string]*domain.Tweet),
	}
}

func (r *MockTweetRepository) Save(tweet *domain.Tweet) error {
	r.tweets[tweet.ID] = tweet
	return nil
}

func (r *MockTweetRepository) FindByUserID(userID string, limit, offset int) ([]*domain.Tweet, error) {
	var result []*domain.Tweet
	for _, tweet := range r.tweets {
		if tweet.UserID == userID {
			result = append(result, tweet)
		}
	}
	return result, nil
}

type MockKafkaWriter struct{}

func (kw *MockKafkaWriter) WriteMessage(message []byte) error {
	return nil
}

func TestCreateTweet(t *testing.T) {
	tweetRepo := NewMockTweetRepository()
	userRepo := NewMockUserRepository()
	kafkaWriter := &MockKafkaWriter{}
	createTweetUseCase := NewCreateTweetUseCase(tweetRepo, userRepo, kafkaWriter)
	user1 := domain.NewUser("user1", "User 1")
	user1.Following["user2"] = true
	tweet, err := createTweetUseCase.Execute("user1", "Hello, world!")
	assert.NoError(t, err)
	assert.NotNil(t, tweet)
	assert.Equal(t, "user1", tweet.UserID)
	assert.Equal(t, "Hello, world!", tweet.Content)
	assert.WithinDuration(t, time.Now(), tweet.Timestamp, time.Second)
}
