package application

import (
	"github.com/google/uuid"
	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/repositories"
)

type CreateTweetUseCase struct {
	tweetRepo repositories.TweetRepository
	userRepo  repositories.UserRepository
	queue     infrastructure.Queue
}

func NewCreateTweetUseCase(tweetRepo repositories.TweetRepository, userRepo repositories.UserRepository, queue infrastructure.Queue) *CreateTweetUseCase {
	return &CreateTweetUseCase{
		userRepo:  userRepo,
		tweetRepo: tweetRepo,
		queue:     queue,
	}
}

func (uc *CreateTweetUseCase) Execute(userID, content string) (*domain.Tweet, error) {
	userExists, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if userExists == nil {
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
