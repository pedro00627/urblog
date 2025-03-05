package application

import (
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/repositories"
)

type FollowUserUseCase interface {
	Execute(followerID, followeeID string) error
}

type followUserUseCase struct {
	userRepo repositories.UserRepository
	queue    infrastructure.Queue
}

func NewFollowUserUseCase(userRepo repositories.UserRepository, queue infrastructure.Queue) FollowUserUseCase {
	return &followUserUseCase{
		userRepo: userRepo,
		queue:    queue,
	}
}

func (uc *followUserUseCase) Execute(followerID, followeeID string) error {
	follower, err := uc.userRepo.FindByID(followerID)
	if err != nil {
		return err
	}
	err = follower.Follow(followeeID)
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
