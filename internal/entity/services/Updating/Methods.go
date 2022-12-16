package Updating

import (
	config "Twitter/config"
	twitterRepository "Twitter/internal/entity/repository"
	twitterService "Twitter/internal/entity/twitterService"
	"fmt"
	"time"
)

type Updating struct {
	twitter    *twitterService.TwitterService
	repository *twitterRepository.TwitterRepository
}

// Create updating
func NewUpdatingStruct(cfg *config.Config) (*Updating, error) {
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

	updating := &Updating{
		twitter:    twitter,
		repository: repository,
	}

	return updating, nil
}

// Func of updating influencers subcriptions(Request + Database method)
func (updating *Updating) UpdateInfluencerSubscriptions(id int) error {
	// Id - id in influencers table, not twitter id

	fmt.Printf("1) Start updating influencer №%v subscriptions\n", id)
	twitterId, err := updating.repository.FindDataWithId(id)
	if err != nil {
		fmt.Printf("Error %v with UpdateInfluencerSubscriptions(twitterId)", err)
		return err
	}

	subscriptions, err := updating.twitter.GetAccountFollowing(twitterId)
	if err != nil {
		fmt.Printf("Error %v with UpdateInfluencerSubscriptions(subscriptions)", err)
		return err
	}

	err = updating.repository.UpdateSubscriptions(id, subscriptions)
	if err != nil {
		fmt.Printf("Error %v with UpdateInfluencerSubscriptions(UpdateSubscriptions(DatabaseMethods))", err)
		return err
	}
	return nil
}

//Main func of updating influencers subcriptions(Infinite for + updating all influencers)
func (updating *Updating) UpdateAllInfluencerSubscriptions() {

	for {

		// Get influencers quantity
		influencerCount, err := updating.repository.GetInfluencersQuantity()
		if err != nil {
			fmt.Printf("Error %v with UpdateAllInfluencerSubscriptions(GetInfluencersQuantity)", err)
		}

		for i := 1; i <= influencerCount; i++ {
			start := time.Now()
			err = updating.UpdateInfluencerSubscriptions(i)
			if err != nil {
				fmt.Printf("Error %v with UpdateAllInfluencerSubscriptions(Updating id: %v)", err, i)
			}

			fmt.Printf("Потратило времени: %v\n", time.Since(start).Seconds())
			time.Sleep(2 * time.Second)
		}

	}
}

// Start UpdateAllInfluencerSubscriptions in goroutine
func (updating *Updating) Start() {
	go updating.UpdateAllInfluencerSubscriptions()
}
