package database

import (
	"antibomberman/mego-post/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func ConnectToDB(cfg *config.Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	maxAttempts := 10
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	for i := 0; i < maxAttempts; i++ {
		log.Printf("connenting to %s", databaseURL)
		db, err = sqlx.Open("postgres", databaseURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Printf("Connected to database")
				return db, nil
			}
		}
		log.Printf("Failed to connect to database, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to database after 10 attempts")
}
