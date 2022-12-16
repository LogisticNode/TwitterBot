package Adding

import (
	config "Twitter/config"
	twitterRepository "Twitter/internal/entity/repository"
	twitterService "Twitter/internal/entity/twitterService"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Adding struct {
	twitter    *twitterService.TwitterService
	repository *twitterRepository.TwitterRepository
}

func NewAddingStruct(cfg *config.Config) (*Adding, error) {
	// Инициализируем twitterService(Requests)
	twitter, err := twitterService.NewTwitterService("", cfg.Api.Key1)
	if err != nil {
		fmt.Printf("Error %v with twitter(NewUpdatingStruct)", err)
	}

	// Инициализируем twitterRepository(DataBase)
	repository, err := twitterRepository.NewRepository(cfg)
	if err != nil {
		fmt.Printf("Error %v with repository(NewUpdatingStruct)", err)
	}

	updating := &Adding{
		twitter:    twitter,
		repository: repository,
	}

	return updating, nil
}

func (adding *Adding) AddInfluencer(username string) error {

	ok, err := adding.repository.ExistsCheck(username)
	if err != nil {
		log.Printf("Error %v with AddInfluencer(ExistsCheck)", err)
		return err
	}

	if ok {
		influencerData, err := adding.twitter.GetAccountData(username)
		if err != nil {
			log.Printf("Error %v with AddInfluencer(influencerData)", err)
			return err
		}

		err = adding.repository.AddInfluencer(influencerData)
		if err != nil {
			log.Printf("Error %v with AddInfluencer(Add to database)", err)
			return err
		}

		log.Printf("%v added successfully", username)
		time.Sleep(2 * time.Second)
		return nil
	} else {
		log.Printf("Username: %v already exists", username)
	}
	return nil
}

func (adding *Adding) AddInfluencers() {
	path := "C:\\Users\\Logistic\\GolandProjects\\Twitter\\internal\\entity\\services\\Adding\\influencers.txt"

	for {
		log.Println("Start reading txt to add influencers")

		readFile, err := os.Open(path)

		if err != nil {
			log.Printf("Error %v with reading file(AddInfluencers)", err)
		}
		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			username := fileScanner.Text()

			err = adding.AddInfluencer(username)
			if err != nil {
				log.Printf("Error %v with AddInfluencer(AddInfluencers)", err)
			}
		}

		readFile.Close()

		err = os.Truncate(path, 0)
		if err != nil {
			log.Printf("Error %v with AddInfluencers(Clearing file)", err)
		}
		time.Sleep(10 * time.Minute)
	}
}

func (adding *Adding) Start() {
	go adding.AddInfluencers()
}
