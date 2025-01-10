package db

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectDB(db_url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	log.Println("Connected to Database!")

	return db, nil
}
