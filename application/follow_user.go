package application

import (
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/db"
)

//go:generate mockgen -destination=./mocks/mock_follow_user.go -package=mocks github.com/pedro00627/urblog/application FollowUser
type FollowUser interface {
	Execute(followerID, followeeID string) error
}

type FollowUserUseCase struct {
	userRepo db.UserRepository
	queue    infrastructure.Queue
}

func NewFollowUserUseCase(userRepo db.UserRepository, queue infrastructure.Queue) FollowUser {
	return &FollowUserUseCase{
		userRepo: userRepo,
		queue:    queue,
	}
}

func (uc *FollowUserUseCase) Execute(followerID, followeeID string) error {
	follower, err := uc.userRepo.FindByID(followerID)
	if err != nil {
		return err
	}
	// find followee
	followee, err := uc.userRepo.FindByID(followeeID)
	if err != nil {
		return err
	}
	err = follower.Follow(followee.ID)
	if err != nil {
		return err
	}
	err = uc.userRepo.Save(follower)
	if err != nil {
		return err
	}
	err = uc.queue.WriteMessage([]byte("User followed: " + followerID))
	if err != nil {
		return err
	}
	return nil
}
