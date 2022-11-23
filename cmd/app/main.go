package main

import (
	"Twitter/internal/entity/twitterService"
	"log"
)

func main() {
	log.Println("Project started")

	// Services
	twitter, err := twitterService.NewTwitterService()
	if err != nil {
		log.Println(err)
	}

	twitter.TestGetAllAccountFollowings()
}
