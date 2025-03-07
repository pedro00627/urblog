package application

import (
	"github.com/google/uuid"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/db"
)

//go:generate mockgen -destination=./mocks/mock_create_tweet.go -package=mocks github.com/pedro00627/urblog/application CreateTweet
type CreateTweet interface {
	Execute(string, string) (*domain.Tweet, error)
}
type CreateTweetUseCase struct {
	tweetRepo db.TweetRepository
	userRepo  db.UserRepository
	queue     infrastructure.Queue
}

func NewCreateTweetUseCase(tweetRepo db.TweetRepository, userRepo db.UserRepository, queue infrastructure.Queue) *CreateTweetUseCase {
	return &CreateTweetUseCase{
		userRepo:  userRepo,
		tweetRepo: tweetRepo,
		queue:     queue,
	}
}

func (uc *CreateTweetUseCase) Execute(userID, content string) (*domain.Tweet, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	tweet, err := domain.NewTweet(generateID(), userID, content)
	if err != nil {
		return nil, err
	}
	err = uc.tweetRepo.Save(tweet)
	if err != nil {
		return nil, err
	}

	err = uc.queue.WriteMessage([]byte("New tweet published: " + tweet.ID))
	if err != nil {
		return nil, err
	}
	return tweet, nil
}

func generateID() string {
	return uuid.New().String()
}
