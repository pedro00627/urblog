package main

import (
	"context"
	"os"
	"time"

	"github.com/pedro00627/urblog/application"
	"github.com/pedro00627/urblog/infrastructure"
	"github.com/pedro00627/urblog/infrastructure/repositories"
	"github.com/pedro00627/urblog/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Dependencies contains the application dependencies
type Dependencies struct {
	TweetController *interfaces.TweetController
	UserController  *interfaces.UserController
}

func InitializeDependencies() (*Dependencies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tweetRepo repositories.TweetRepository
	var userRepo repositories.UserRepository
	var queue infrastructure.Queue

	//Creating Repositories
	if os.Getenv("DATABASE") == "" {
		tweetRepo = repositories.NewInMemoryTweetRepository()
		userRepo = repositories.NewInMemoryUserRepository()
	} else {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
		if err != nil {
			return nil, err
		}
		db := client.Database(os.Getenv("DATABASE"))
		tweetRepo = repositories.NewMongoTweetRepository(db)
		userRepo = repositories.NewMongoUserRepository(db)
	}

	if kafkaBroker := os.Getenv("KAFKA_BROKER"); kafkaBroker != "" {
		queue = infrastructure.NewKafkaWriter(kafkaBroker)
	} else {
		queue = infrastructure.NewInMemoryqueue()
	}

	// Creating Use Cases
	createTweet := application.NewCreateTweetUseCase(tweetRepo, userRepo, queue)
	followUser := application.NewFollowUserUseCase(userRepo, queue)
	getTimeline := application.NewGetTimelineUseCase(tweetRepo, userRepo)
	loadUsersUseCase := application.NewLoadUsersUseCase(userRepo)

	// Creating Controllers
	tweetController := interfaces.NewTweetController(createTweet)
	userController := interfaces.NewUserController(followUser, getTimeline, loadUsersUseCase)

	deps := &Dependencies{
		TweetController: tweetController,
		UserController:  userController,
	}

	return deps, nil
}
