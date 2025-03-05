package repositories

import "github.com/pedro00627/urblog/domain"

type TweetRepository interface {
	FindByUserID(string, int, int) ([]*domain.Tweet, error)
	Save(*domain.Tweet) error
}

type UserRepository interface {
	FindByID(string) (*domain.User, error)
	FindByName(s string) (*domain.User, error)
	Save(*domain.User) error
}
