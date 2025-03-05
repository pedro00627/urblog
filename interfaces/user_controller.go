package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/pedro00627/urblog/application"
)

type UserController struct {
	followUserUseCase  application.FollowUserUseCase
	getTimelineUseCase application.GetTimelineUseCase
	loadUsersUseCase   *application.LoadUsersUseCase
}

func NewUserController(followUserUseCase application.FollowUserUseCase, getTimelineUseCase application.GetTimelineUseCase, loadUsersUseCase *application.LoadUsersUseCase) *UserController {
	return &UserController{
		followUserUseCase:  followUserUseCase,
		getTimelineUseCase: getTimelineUseCase,
		loadUsersUseCase:   loadUsersUseCase,
	}
}

func (c *UserController) FollowUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FollowerID string `json:"follower_id"`
		FolloweeID string `json:"followee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := c.followUserUseCase.Execute(req.FollowerID, req.FolloweeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) GetTimeline(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Limit  int    `json:"limit,omitempty"`
		Offset int    `json:"offset,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tweets, err := c.getTimelineUseCase.Execute(req.UserID, req.Limit, req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := make([]struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}, len(tweets))
	for i, tweet := range tweets {
		resp[i] = struct {
			ID        string `json:"id"`
			UserID    string `json:"user_id"`
			Content   string `json:"content"`
			Timestamp string `json:"timestamp"`
		}{
			ID:        tweet.ID,
			UserID:    tweet.UserID,
			Content:   tweet.Content,
			Timestamp: tweet.Timestamp.String(),
		}
	}
	w.Header().Set("Content-Type", "application/json")
_:
	json.NewEncoder(w).Encode(resp)
}

func (c *UserController) LoadUsers(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "file parameter is required", http.StatusBadRequest)
		return
	}

	users, err := c.loadUsersUseCase.Execute(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
_:
	json.NewEncoder(w).Encode(users)
}
