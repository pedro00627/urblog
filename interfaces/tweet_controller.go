package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/pedro00627/urblog/application"
)

type TweetController struct {
	createTweet application.CreateTweet
}

func NewTweetController(createTweet application.CreateTweet) *TweetController {
	return &TweetController{
		createTweet: createTweet,
	}
}

// @Summary Create a new tweet
// @Description Create a new tweet with the given content
// @Tags tweets
// @Accept  json
// @Produce  json
// @Param   tweet  body  Tweet  true  "Tweet content"
// @Success 200 {object} Tweet
// @Failure 400 {object} ErrorResponse
// @Router /tweets [post]
func (c *TweetController) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tweet, err := c.createTweet.Execute(req.UserID, req.Content)
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
