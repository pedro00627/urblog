package domain

import "time"

type Tweet struct {
	ID        string
	UserID    string
	Content   string
	Timestamp time.Time
}

func NewTweet(id, userID, content string) (*Tweet, error) {
	if len(content) == 0 || len(content) > 280 {
		return nil, ErrInvalidTweetContent
	}
	return &Tweet{
		ID:        id,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}, nil
}
