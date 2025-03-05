package application

import (
	"log"
	"sort"

	"github.com/pedro00627/urblog/domain"
	"github.com/pedro00627/urblog/infrastructure/repositories"
)

type GetTimelineUseCase interface {
	Execute(userID string, limit, offset int) ([]*domain.Tweet, error)
}

type getTimelineUseCase struct {
	tweetRepo repositories.TweetRepository
	userRepo  repositories.UserRepository
}

func NewGetTimelineUseCase(tweetRepo repositories.TweetRepository, userRepo repositories.UserRepository) GetTimelineUseCase {
	return &getTimelineUseCase{
		tweetRepo: tweetRepo,
		userRepo:  userRepo,
	}
}

func (uc *getTimelineUseCase) Execute(userID string, limit, offset int) ([]*domain.Tweet, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var allTweets []*domain.Tweet
	for followedUserID := range user.Following {
		// get user name
		followedUserId, err := uc.userRepo.FindByName(followedUserID)
		if err != nil {
			log.Printf("Error finding user by name: %s", followedUserID)
			continue
		}

		tweets, err := uc.tweetRepo.FindByUserID(followedUserId.ID, limit, offset)
		if err != nil {
			return nil, err
		}
		allTweets = append(allTweets, tweets...)
	}

	// Ordenar los tweets por fecha
	sort.Slice(allTweets, func(i, j int) bool {
		return allTweets[i].Timestamp.After(allTweets[j].Timestamp)
	})

	// Aplicar paginación
	start := offset
	end := offset + limit
	if start > len(allTweets) {
		start = len(allTweets)
	}
	if end > len(allTweets) {
		end = len(allTweets)
	}
	paginatedTweets := allTweets[start:end]

	return paginatedTweets, nil
}
