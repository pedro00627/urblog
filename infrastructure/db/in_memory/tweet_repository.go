package in_memory

import "github.com/pedro00627/urblog/domain"

type InMemoryTweetRepository struct {
	tweets         map[string]*domain.Tweet
	tweetsByUserID map[string][]string
}

func NewInMemoryTweetRepository() *InMemoryTweetRepository {
	return &InMemoryTweetRepository{
		tweets:         make(map[string]*domain.Tweet),
		tweetsByUserID: make(map[string][]string),
	}
}

func (r *InMemoryTweetRepository) Save(tweet *domain.Tweet) error {
	r.tweets[tweet.ID] = tweet
	r.tweetsByUserID[tweet.UserID] = append(r.tweetsByUserID[tweet.UserID], tweet.ID)
	return nil
}

func (r *InMemoryTweetRepository) FindByUserID(userID string, limit, offset int) ([]*domain.Tweet, error) {
	var result []*domain.Tweet
	count := 0

	for _, tweet := range r.tweets {
		if tweet.UserID == userID {
			if count >= offset && count < offset+limit {
				result = append(result, tweet)
			}
			count++
		}
	}

	return result, nil
}
