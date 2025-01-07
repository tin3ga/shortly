package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/thanhpk/randstr"

	"github.com/tin3ga/shortly/internal/database"

	_ "github.com/lib/pq"
	_ "github.com/tin3ga/shortly/docs"
)

const version = "0.1.1"

// Shorten Link model info
//
//	@Description	Shorten link Model
//	@Description	Url, Custom_alias
type ShortenLink struct {
	Url          string `json:"url"`
	Custom_alias string `json:"custom_alias"`
	// expiration_date string `json:"expiration_date"`

}

// Delete Link model info
//
//	@Description	Delete Link Model
//	@Description	Url
type DeleteLink struct {
	Url string `json:"url"`
}

// Verify Connection
//
//	@Summary		Checks connectivity
//	@Description	Returns pong
//	@Success		200
//	@Router			/ [get]
func ping(c *fiber.Ctx) error {
	log.Println("Ping.....")
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}

// getLink Fetch a Original URL by Short URL
//
//	@Summary		Fetch a Original URL by Short URL
//	@Description	Redirects to the original URL
//	@Param			link	path	string	true	"Redirects to Original URL"
//	@Success		301
//	@Failure		404
//	@Router			/api/v1/{link} [get]
func getLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	link := c.Params("link")

	data, err := queries.GetLongLink(ctx, link)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found"})
	}
	log.Println("Redirecting to: ", data.LongLink)
	return c.Redirect(data.LongLink, fiber.StatusMovedPermanently)

}

// shortenLink Insert an entry for a Short URL and Long URL
//
//	@Summary		Insert an entry for a Short URL and Long URL
//	@Description	Returns a Short URL
//	@Param			shorten_link	body	ShortenLink	true	"Shorten a Link (custom alias is optional)"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/api/v1/shorten [post]
func shortenLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context, urlStr string) error {
	url := new(ShortenLink)

	if err := c.BodyParser(url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var ShortLink string

	if url.Custom_alias == "" {
		ShortLink = randstr.Hex(8) // Generate a random 8 character string
	} else {
		ShortLink = url.Custom_alias
	}

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

// deleteLink Delete url data by short url
//
//	@Summary		Delete url data by short url
//	@Description	Returns a success message
//	@Param			url	body	DeleteLink	true	"Delete a Link"
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/shorten [delete]
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

// getLinks Fetch all links
//
//	@Summary		Fetch all links
//	@Description	Returns all links
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/ [get]
func getLinks(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	data, err := queries.GetLinks(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cannot fetch links"})
	}
	log.Println("Fetching links")
	return c.JSON(data)

}

//	@title			Shortly API
//	@version		0.1.1
//	@description	This is a URL Shortener backend API built with Go.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.email	tinegagideon@outlook.com
//	@license.name	MIT License
//	@license.url	https://mit-license.org/
//	@host			shortly-5p7d.onrender.com
//	@BasePath		/

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	connStr := os.Getenv("DATABASE_URL")
	urlStr := os.Getenv("URL")

	log.Printf("Starting server on port %v", portString)
	log.Printf("Serving version %v", version)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database!")

	defer db.Close()

	ctx := context.Background()
	queries := database.New(db)

	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/", ping)
	app.Get("/api/v1/:link", func(c *fiber.Ctx) error {
		return getLink(c, queries, ctx)
	})
	app.Post("/api/v1/shorten", func(c *fiber.Ctx) error {
		return shortenLink(c, queries, ctx, urlStr)
	})
	app.Delete("/api/v1/shorten", func(c *fiber.Ctx) error {
		return deleteLink(c, queries, ctx)
	})
	app.Get("/api/v1/", func(c *fiber.Ctx) error {
		return getLinks(c, queries, ctx)
	})

	app.Listen(":" + portString)
}
