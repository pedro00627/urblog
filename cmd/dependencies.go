package main

import (
	"context"
	"github.com/pedro00627/urblog/infrastructure/db"
	"github.com/pedro00627/urblog/infrastructure/db/in_memory"
	mongo2 "github.com/pedro00627/urblog/infrastructure/db/mongo"
	inmemory2 "github.com/pedro00627/urblog/infrastructure/queue/in_memory"
	"github.com/pedro00627/urblog/infrastructure/queue/kafka"
	"os"
	"time"

	"github.com/pedro00627/urblog/application"
	"github.com/pedro00627/urblog/infrastructure"
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

	var tweetRepo db.TweetRepository
	var userRepo db.UserRepository
	var queue infrastructure.Queue

	//Creating Repositories
	if os.Getenv("DATABASE") == "" {
		tweetRepo = in_memory.NewInMemoryTweetRepository()
		userRepo = in_memory.NewInMemoryUserRepository()
	} else {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
		if err != nil {
			return nil, err
		}
		database := client.Database(os.Getenv("DATABASE"))
		tweetRepo = mongo2.NewTweetRepository(database)
		userRepo = mongo2.NewUserRepository(database)
	}

	if kafkaBroker := os.Getenv("KAFKA_BROKER"); kafkaBroker != "" {
		queue = kafka.NewWriter(kafkaBroker)
	} else {
		queue = inmemory2.NewInMemoryQueue()
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
