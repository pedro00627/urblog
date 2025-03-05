package application

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

type MockKafkaWriter struct{}

func (kw *MockKafkaWriter) WriteMessage(message []byte) error {
	return nil
}

func TestCreateTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserRepository(ctrl)
	tweetRepo := mocks.NewMockTweetRepository(ctrl)

	userRepo.EXPECT().FindByID("user1").Return(domain.NewUser("user1", "User 1"), nil).Times(1)
	tweetRepo.EXPECT().Save(gomock.Any()).Times(1)

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
