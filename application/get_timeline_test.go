package application

import (
	"testing"
	"time"

	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeline(t *testing.T) {
	tweetRepo := NewMockTweetRepository()
	userRepo := NewMockUserRepository()
	getTimelineUseCase := NewGetTimelineUseCase(tweetRepo, userRepo)

	user1 := domain.NewUser("user1", "User One")
	user2 := domain.NewUser("user2", "User Two")
	user1.Follow("user2")
	userRepo.Save(user1)
	userRepo.Save(user2)

	tweetRepo.Save(&domain.Tweet{
		ID:        "tweet1",
		UserID:    "user2",
		Content:   "Tweet from user2",
		Timestamp: time.Now(),
	})

	timeline, err := getTimelineUseCase.Execute("user1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, timeline, 1)
	assert.Equal(t, "Tweet from user2", timeline[0].Content)
}
