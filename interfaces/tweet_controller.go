package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/pedro00627/urblog/application"
)

type TweetController struct {
	createTweetUseCase *application.CreateTweetUseCase
}

func NewTweetController(createTweetUseCase *application.CreateTweetUseCase) *TweetController {
	return &TweetController{
		createTweetUseCase: createTweetUseCase,
	}
}

func (c *TweetController) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tweet, err := c.createTweetUseCase.Execute(req.UserID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := struct {
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
	w.Header().Set("Content-Type", "application/json")
_:
	json.NewEncoder(w).Encode(resp)
}
