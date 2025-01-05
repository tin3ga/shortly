package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/thanhpk/randstr"

	"github.com/tin3ga/shortly/internal/database"

	_ "github.com/lib/pq"
)

type ShortenLink struct {
	Url string `json:"url"`
	// custom_alias string `json:"custom_alias"`
	// expiration_date string `json:"expiration_date"`

}
type DeleteLink struct {
	Url string `json:"url"`
}

func ping(c *fiber.Ctx) error {
	log.Println("Ping.....")
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
func getLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	link := c.Params("link")

	data, err := queries.GetLongLink(ctx, link)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found"})
	}
	log.Println("Redirecting to: ", data.LongLink)
	return c.Redirect(data.LongLink, fiber.StatusMovedPermanently)

}

func shortenLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context, urlStr string) error {
	url := new(ShortenLink)

	if err := c.BodyParser(url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	ShortLink := randstr.Hex(8) // Generate a random 8 character string
	LongLink := url.Url
	uuid := uuid.New()

	params := database.CreateShortLinkParams{
		ID:        uuid,
		ShortLink: ShortLink,
		LongLink:  LongLink,
	}
	_, err := queries.CreateShortLink(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create short link"})
	}
	log.Println("Created a shortend link: ", ShortLink)

	return c.JSON(fiber.Map{"Success": "Shortened link created", "url": urlStr + ShortLink})
}

func deleteLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	url := new(DeleteLink)

	if err := c.BodyParser(url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check if the short link exists
	_, err := queries.GetLongLink(ctx, url.Url)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found"})
	}

	err = queries.DeleteLink(ctx, url.Url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete short link"})
	}

	log.Println("Deleted a shortend link: ", url.Url)
	return c.JSON(fiber.Map{"Success": "Shortened link deleted"})

}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")
	urlStr := os.Getenv("URL")

	log.Printf("Starting server on port %v", portString)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database!")

	defer db.Close()

	ctx := context.Background()
	queries := database.New(db)

	// data, err := queries.ListLinks(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, item := range data {
	// 	fmt.Println(item.ShortLink, item.LongLink)
	// }

	app := fiber.New()

	app.Get("", ping)
	app.Get("/api/v1/:link", func(c *fiber.Ctx) error {
		return getLink(c, queries, ctx)
	})
	app.Post("/api/v1/shorten", func(c *fiber.Ctx) error {
		return shortenLink(c, queries, ctx, urlStr)
	})
	app.Delete("/api/v1/shorten", func(c *fiber.Ctx) error {
		return deleteLink(c, queries, ctx)
	})

	app.Listen(":" + portString)
}
