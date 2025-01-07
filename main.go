package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"

	"github.com/tin3ga/shortly/internal/database"

	_ "github.com/lib/pq"
	_ "github.com/tin3ga/shortly/docs"
)

const version = "0.2.0"

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
func getLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context, rdb *redis.Client, ttl time.Duration) error {
	link := c.Params("link")

	// caching - Get

	// check if rdb is not nil to prevent nil pointer dereference error
	if rdb != nil {
		val, err := rdb.Get(ctx, link).Result()
		if err != nil {
			log.Printf("Cannot find data with key: %s", link)
		}

		if val != "" {

			var data database.Shortly
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				log.Printf("Error unmarshaling data: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
			}

			log.Println("Redirecting to: ", data.LongLink)
			return c.Redirect(data.LongLink, fiber.StatusMovedPermanently)

		}

	}
	// end caching - Get

	// check if value in database, returns if no data is found skips caching set
	data, err := queries.GetLongLink(ctx, link)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found"})
	}

	// caching - Set
	// cache results if data exists key is short url/link
	if rdb != nil {
		marshalData, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		}
		// add a item to cache if it did not exist

		err = rdb.Set(ctx, link, marshalData, ttl).Err()
		if err != nil {
			log.Print(err)
		} else {
			log.Print("Item added to cache")
		}

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
	log.Println("Created a shortened link: ", ShortLink)

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

	log.Println("Deleted a shortened link: ", url.Url)
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
//	@version		0.2.0
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
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	enableCaching := os.Getenv("cachingEnabled")
	ttlStr := os.Getenv("cache_ttl")

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

	// Redis setup
	var cachingEnabled bool
	if enableCaching != "" {
		cachingEnabled, err = strconv.ParseBool(enableCaching)
		if err != nil {
			log.Fatalf("Invalid ENABLE_CACHING value: %v", err)
		}
	}
	log.Printf("Caching Enabled: %v", cachingEnabled)

	// Convert Redis cache ttl from string to int
	var ttlInt int
	if ttlStr != "" {
		var err error
		ttlInt, err = strconv.Atoi(redisDB)
		if err != nil {
			log.Fatalf("Invalid cache_ttl: %v", err)
		}
	}

	ttl := time.Duration(ttlInt) * time.Minute

	var rdb *redis.Client

	if cachingEnabled {
		rdb, err = initializeRedis(ctx, redisAddr, redisPassword, redisDB)

		if err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}
		defer rdb.Close() // Close Redis client when no longer needed

	}

	log.Printf("Redis Connection: %v", rdb)

	// end Redis setup

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
		return getLink(c, queries, ctx, rdb, ttl)
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

func initializeRedis(ctx context.Context, redisAddr string, redisPassword string, redisDB string) (*redis.Client, error) {
	// Convert Redis DB from string to int
	var redisDBInt int
	if redisDB != "" {
		var err error
		redisDBInt, err = strconv.Atoi(redisDB)
		if err != nil {
			log.Fatalf("Invalid REDIS_DB: %v", err)
		}
	}
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDBInt,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Cannot connect to redis db: \nError %v", err)
		return nil, err
	}

	if pong == "PONG" {
		log.Print("Connected to cache server")

	}

	return rdb, nil
}
