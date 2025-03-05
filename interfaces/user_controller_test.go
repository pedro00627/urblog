package interfaces

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
)

type MockFollowUserUseCase struct {
	err error
}

func (uc *MockFollowUserUseCase) Execute(followerID, followeeID string) error {
	return uc.err
}

type MockGetTimelineUseCase struct {
	tweets []*domain.Tweet
	err    error
}

func (uc *MockGetTimelineUseCase) Execute(userID string, limit, offset int) ([]*domain.Tweet, error) {
	return uc.tweets, uc.err
}

func TestFollowUser(t *testing.T) {
	mockFollowUserUseCase := &MockFollowUserUseCase{}
	mockGetTimelineUseCase := &MockGetTimelineUseCase{}

	userController := NewUserController(mockFollowUserUseCase, mockGetTimelineUseCase, nil)

	reqBody := bytes.NewBufferString(`{"follower_id": "user1", "followee_id": "user2"}`)
	req := httptest.NewRequest(http.MethodPost, "/follow", reqBody)
	w := httptest.NewRecorder()

	userController.FollowUser(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGetTimeline(t *testing.T) {
	tweet := &domain.Tweet{
		ID:        "tweet1",
		UserID:    "user2",
		Content:   "Tweet from user2",
		Timestamp: time.Now(),
	}
	mockGetTimelineUseCase := &MockGetTimelineUseCase{
		tweets: []*domain.Tweet{tweet},
	}
	mockFollowUserUseCase := &MockFollowUserUseCase{}
	userController := NewUserController(mockFollowUserUseCase, mockGetTimelineUseCase, nil)

	reqBody := bytes.NewBufferString(`{"user_id": "user1", "limit": 10, "offset": 0}`)
	req := httptest.NewRequest(http.MethodPost, "/timeline", reqBody)
	w := httptest.NewRecorder()

	userController.GetTimeline(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			
		}
	}(resp.Body)

	var timeline []struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}
_:
	json.NewDecoder(resp.Body).Decode(&timeline)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, timeline, 1)
	assert.Equal(t, "tweet1", timeline[0].ID)
	assert.Equal(t, "user2", timeline[0].UserID)
	assert.Equal(t, "Tweet from user2", timeline[0].Content)
}
