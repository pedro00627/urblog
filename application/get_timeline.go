package application

import (
	"github.com/pedro00627/urblog/infrastructure/db"
	"log"
	"sort"

	"github.com/pedro00627/urblog/domain"
)

//go:generate mockgen -destination=./mocks/mock_get_timeline.go -package=mocks github.com/pedro00627/urblog/application GetTimeline
type GetTimeline interface {
	Execute(userID string, limit, offset int) ([]*domain.Tweet, error)
}

type GetTimelineUseCase struct {
	tweetRepo db.TweetRepository
	userRepo  db.UserRepository
}

func NewGetTimelineUseCase(tweetRepo db.TweetRepository, userRepo db.UserRepository) GetTimeline {
	return &GetTimelineUseCase{
		tweetRepo: tweetRepo,
		userRepo:  userRepo,
	}
}

func (uc *GetTimelineUseCase) Execute(userID string, limit, offset int) ([]*domain.Tweet, error) {
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

	// Sort tweets by timestamp
	sortTweetsByNewest(allTweets)

	// Apply pagination
	paginatedTweets := paginateTweets(offset, limit, allTweets)

	return paginatedTweets, nil
}

func paginateTweets(offset int, limit int, allTweets []*domain.Tweet) []*domain.Tweet {
	start := offset
	end := offset + limit
	if start > len(allTweets) {
		start = len(allTweets)
	}
	if end > len(allTweets) {
		end = len(allTweets)
	}
	paginatedTweets := allTweets[start:end]
	return paginatedTweets
}

func sortTweetsByNewest(allTweets []*domain.Tweet) {
	sort.Slice(allTweets, func(i, j int) bool {
		return allTweets[i].Timestamp.After(allTweets[j].Timestamp)
	})
}
