package db

import "github.com/pedro00627/urblog/domain"

//go:generate mockgen -destination=./mocks/mock_tweet_repository.go -package=mocks github.com/pedro00627/urblog/infrastructure/repositories TweetRepository
//go:generate mockgen -destination=./mocks/mock_user_repository.go -package=mocks github.com/pedro00627/urblog/infrastructure/repositories UserRepository

type TweetRepository interface {
	FindByUserID(string, int, int) ([]*domain.Tweet, error)
	Save(*domain.Tweet) error
}

type UserRepository interface {
	FindByID(string) (*domain.User, error)
	FindByName(s string) (*domain.User, error)
	Save(*domain.User) error
}
