package application

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/repositories/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTweetRepo := mocks.NewMockTweetRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	userID := "user1"
	followedUserID := "user2"
	tweet1 := &domain.Tweet{
		ID:        "tweet1",
		UserID:    followedUserID,
		Content:   "Hello, world!",
		Timestamp: time.Now().Add(-1 * time.Hour),
	}
	tweet2 := &domain.Tweet{
		ID:        "tweet2",
		UserID:    followedUserID,
		Content:   "Another tweet",
		Timestamp: time.Now().Add(-2 * time.Hour),
	}

	user := &domain.User{
		ID:       userID,
		Username: "user1",
		Following: map[string]bool{
			followedUserID: true,
		},
	}

	mockUserRepo.EXPECT().FindByID(userID).Return(user, nil)
	mockUserRepo.EXPECT().FindByName(followedUserID).Return(&domain.User{ID: followedUserID}, nil)
	mockTweetRepo.EXPECT().FindByUserID(followedUserID, 10, 0).Return([]*domain.Tweet{tweet1, tweet2}, nil)

	useCase := NewGetTimelineUseCase(mockTweetRepo, mockUserRepo)
	tweets, err := useCase.Execute(userID, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, tweets, 2)
	assert.Equal(t, tweet1.ID, tweets[0].ID)
	assert.Equal(t, tweet2.ID, tweets[1].ID)
}
