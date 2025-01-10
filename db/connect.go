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

	log.Println("Connected to Database!")

	return db, nil
}
