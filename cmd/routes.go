package main

import (
	"net/http"
)

func ConfigureRoutes(mux *http.ServeMux, deps *Dependencies) {
	mux.HandleFunc("/tweets", deps.TweetController.CreateTweet)
	mux.HandleFunc("/follow", deps.UserController.FollowUser)
	mux.HandleFunc("/timeline", deps.UserController.GetTimeline)
	mux.HandleFunc("/load-users", deps.UserController.LoadUsers)
}
