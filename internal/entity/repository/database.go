package repository

import (
	"Twitter/config"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type TwitterRepository struct {
	db *sql.DB
}

// Create repository
func NewRepository(cfg *config.Config) (*TwitterRepository, error) {

	path := fmt.Sprintf("%s:%s@/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.DbName)

	// Connecting to database
	db, err := sql.Open("mysql", path)

	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	// Dont know
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Minute * 90)
	log.Printf("Connected to DB %s successfully\n", cfg.Database.DbName)

	repository := &TwitterRepository{
		db: db,
	}
	return repository, nil

}
