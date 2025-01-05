package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/tin3ga/shortly/internal/database"

	_ "github.com/lib/pq"
)

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")

	log.Printf("Starting server on port %v", portString)
	log.Printf("Connecting to database")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database!")

	defer db.Close()

	ctx := context.Background()
	queries := database.New(db)

	data, err := queries.ListLinks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range data {
		fmt.Println(item.ShortLink, item.LongLink)
	}

	app := fiber.New()

	app.Get("/", ping)

	app.Listen(":" + portString)
}
