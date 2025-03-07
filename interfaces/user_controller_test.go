package interfaces

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/pedro00627/urblog/application/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedro00627/urblog/domain"
	"github.com/stretchr/testify/assert"
)

func TestFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFollowUser := mocks.NewMockFollowUser(ctrl)
	mockGetTimeline := mocks.NewMockGetTimeline(ctrl)
	mockLoadUsers := mocks.NewMockLoadUsers(ctrl)

	userController := NewUserController(mockFollowUser, mockGetTimeline, mockLoadUsers)

	t.Run("success", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"follower_id": "user1", "followee_id": "user2"}`)
		req := httptest.NewRequest(http.MethodPost, "/follow", reqBody)
		w := httptest.NewRecorder()

		mockFollowUser.EXPECT().Execute("user1", "user2").Return(nil).Times(1)

		userController.FollowUser(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("error", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"follower_id": "user1", "followee_id": "user2"}`)
		req := httptest.NewRequest(http.MethodPost, "/follow", reqBody)
		w := httptest.NewRecorder()

		mockFollowUser.EXPECT().Execute("user1", "user2").Return(assert.AnError).Times(1)

		userController.FollowUser(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetTimelineUseCase := mocks.NewMockGetTimeline(ctrl)
	mockFollowUserUseCase := mocks.NewMockFollowUser(ctrl)
	mockLoadUsers := mocks.NewMockLoadUsers(ctrl)

	userController := NewUserController(mockFollowUserUseCase, mockGetTimelineUseCase, mockLoadUsers)

	tweet := &domain.Tweet{
		ID:        "tweet1",
		UserID:    "user2",
		Content:   "Tweet from user2",
		Timestamp: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"user_id": "user1", "limit": 10, "offset": 0}`)
		req := httptest.NewRequest(http.MethodPost, "/timeline", reqBody)
		w := httptest.NewRecorder()

		mockGetTimelineUseCase.EXPECT().Execute("user1", 10, 0).Return([]*domain.Tweet{tweet}, nil).Times(1)

		userController.GetTimeline(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		var timeline []struct {
			ID        string `json:"id"`
			UserID    string `json:"user_id"`
			Content   string `json:"content"`
			Timestamp string `json:"timestamp"`
		}
		json.NewDecoder(resp.Body).Decode(&timeline)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, timeline, 1)
		assert.Equal(t, "tweet1", timeline[0].ID)
		assert.Equal(t, "user2", timeline[0].UserID)
		assert.Equal(t, "Tweet from user2", timeline[0].Content)
	})

	t.Run("error", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"user_id": "user1", "limit": 10, "offset": 0}`)
		req := httptest.NewRequest(http.MethodPost, "/timeline", reqBody)
		w := httptest.NewRecorder()

		mockGetTimelineUseCase.EXPECT().Execute("user1", 10, 0).Return(nil, assert.AnError).Times(1)

		userController.GetTimeline(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestLoadUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFollowUser := mocks.NewMockFollowUser(ctrl)
	mockGetTimeline := mocks.NewMockGetTimeline(ctrl)
	mockLoadUsers := mocks.NewMockLoadUsers(ctrl)

	userController := NewUserController(mockFollowUser, mockGetTimeline, mockLoadUsers)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/load_users?file=users.json", nil)
		w := httptest.NewRecorder()

		mockLoadUsers.EXPECT().Execute("users.json").Return([]domain.User{
			{ID: "user1", Username: "User One"},
			{ID: "user2", Username: "User Two"},
		}, nil).Times(1)

		userController.LoadUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		var users []domain.User
		json.NewDecoder(resp.Body).Decode(&users)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Len(t, users, 2)
		assert.Equal(t, "user1", users[0].ID)
		assert.Equal(t, "User One", users[0].Username)
		assert.Equal(t, "user2", users[1].ID)
		assert.Equal(t, "User Two", users[1].Username)
	})

	t.Run("error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/load_users?file=users.json", nil)
		w := httptest.NewRecorder()

		mockLoadUsers.EXPECT().Execute("users.json").Return(nil, assert.AnError).Times(1)

		userController.LoadUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
