package db

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Cannot connect to postgres db: \nError %v", err)
		return nil, err
	}

	log.Print("Connected to Postgres DB")

	return db, nil
}
