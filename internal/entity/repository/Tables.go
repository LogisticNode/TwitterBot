package repository

import (
	"context"
	"log"
	"time"
)

//   Creation tables   // v

// Table of influencers
func (repository *TwitterRepository) CreateInfluencersTable() error {
	query := `CREATE TABLE IF NOT EXISTS influencers(
	id int primary key auto_increment, 
	name text,
	username text NOT NULL, 
	twitter_id text NOT NULL,
	followers int(15) NOT NULL,
	score int,
	created_at datetime default CURRENT_TIMESTAMP,
	updated_at datetime default CURRENT_TIMESTAMP)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := repository.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating twitters table\n", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected\n", err)
		return err
	}
	log.Println("Influencers table created(if doesnt exists)")
	return nil

}

// Table of subscriptions
func (repository *TwitterRepository) CreateSubscriptionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS subscriptions(
	id int primary key auto_increment,
	influencer_id int NOT NULL,
	foreign key (influencer_id) references influencers (id),
	subscription_id text NOT NULL,
	updated_at datetime default CURRENT_TIMESTAMP)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := repository.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating subscriptions table\n", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected\n", err)
		return err
	}
	log.Println("Subscriptions table created(if doesnt exists)")
	return nil
}
