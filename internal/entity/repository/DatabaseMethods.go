package repository

import (
	config "Twitter/config"
	twitter "Twitter/internal/entity/twitterService"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TwitterRepository struct {
	db *sql.DB
}

// Create repository
func NewRepository(cfg *config.Config) (*TwitterRepository, error) {

	path := fmt.Sprintf("%s:%s@/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.DbName)

	//Подключение к базе данных
	db, err := sql.Open("mysql", path)

	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	//Надо подумать
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 90)
	log.Printf("Connected to DB %s successfully\n", cfg.Database.DbName)

	repository := &TwitterRepository{
		db: db,
	}
	return repository, nil

}

// Add influencer to database
func (repository *TwitterRepository) AddInfluencer(influencer twitter.AccountData) error {
	query := "INSERT INTO influencers(name, username, twitter_id, followers) VALUES (?, ?, ?, ?)"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement\n", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, influencer.Name, influencer.Username, influencer.User_id, influencer.Follower_count)
	if err != nil {
		log.Printf("Error %s when inserting row into influencers table(AddInfluencer)\n", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected(AddInfluencer)\n", err)
		return err
	}

	if rows != 1 {
		log.Printf("Error %s with AddInfluencer() rows\n", err)
		return err
	}
	return nil
}

//   Subscriptions   //

// Add influencer subscriptions to database
func (repository *TwitterRepository) AddInfluencerSubscriptions(id int, subscriptions twitter.AccountFollowing) error {
	query := "INSERT INTO subscriptions(influencer_id, subscription_id) VALUES (?, ?)"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	for i := 0; i <= len(subscriptions.Result)-1; i++ {
		res, err := stmt.ExecContext(ctx, id, subscriptions.Result[i].User_id)
		if err != nil {
			log.Printf("Error %s when inserting row into subscriptions table\n", err)
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			log.Printf("Error %s when finding rows affected(AddInfluencerSubscription)\n", err)
		}

		if rows != 1 {
			log.Printf("Error %s with AddInfluencerSubscription() rows\n", err)
			return err
		}

	}

	fmt.Printf("  3)%v rows added to Influencer(%v)\n", len(subscriptions.Result), id)
	return nil
}

// Delete influencer subscriptions(to update)
func (repository *TwitterRepository) DeleteInfluencerSubscriptions(id int) error {
	query := "DELETE FROM subscriptions WHERE influencer_id = ?"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Printf("Error %s when deleting row in subscriptions table\n", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected(DeleteInfluencerFollowed)\n", err)
		return err
	}

	fmt.Printf("  2)%v rows is deleted to Influencer(%v)\n", rows, id)
	return nil
}

// DeleteInfluencerSubscriptions + AddInfluencerSubscriptions
func (repository *TwitterRepository) UpdateSubscriptions(id int, subscriptions twitter.AccountFollowing) error {

	err := repository.DeleteInfluencerSubscriptions(id)
	if err != nil {
		fmt.Println("Error with UpdateSubscriptions(Deleting)\n")
		return err
	}

	err = repository.AddInfluencerSubscriptions(id, subscriptions)
	if err != nil {
		fmt.Println("Error with UpdateSubscriptions(Adding)\n")
		return err
	}

	fmt.Printf("4)Subscriptions on influencer(%v) updated\n", id)
	return nil
}

// Dont work, this is to discord bot
func (repository *TwitterRepository) FindSubscriptionsWithId(influencerId int) ([]string, error) {

	var subscriptions []string
	var subscription string

	rows, err := repository.db.Query("SELECT subscription_id FROM subscriptions WHERE influencer_id = ?", influencerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&subscription)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// Get influencer twitter id to update his data
func (repository *TwitterRepository) FindDataWithId(id int) (string, error) {
	var twitterId string

	rows, err := repository.db.Query("SELECT twitter_id FROM influencers WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&twitterId)
		if err != nil {
			fmt.Printf("Error %v with FindDataWithId", err)
			return "", err
		}
		return twitterId, nil
	}
	return "", nil
}

// Get all influencers count to evaluate time(All updating time = 1 hour)
func (repository *TwitterRepository) GetInfluencersQuantity() (int, error) {
	var count int

	rows, err := repository.db.Query("SELECT COUNT(id) FROM influencers")
	if err != nil {
		return 1, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Printf("Error %v with FindDataWithId", err)
			return 1, err
		}
		return count, nil
	}
	return 1, nil

}

// Check influencer exists to Adding(AddInfluencer)
func (repository *TwitterRepository) ExistsCheck(username string) (bool, error) {
	var count int

	rows, err := repository.db.Query("SELECT COUNT(id) FROM influencers WHERE username = ?", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Printf("Error %v with FindDataWithId", err)
		}
		if count == 0 {
			return true, err
		}
		return false, err
	}
	return false, nil
}

// To fing new subscriptions
func (repository *TwitterRepository) FindMismatches(influencerId int) error {

	subscriptions, err := repository.FindSubscriptionsWithId(influencerId)
	if err != nil {
		log.Printf("Error %v with FindMismatches(FindSubscriptionsWithId)", err)
		return err
	}

}
