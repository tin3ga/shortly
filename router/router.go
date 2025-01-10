package router

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/tin3ga/shortly/handler"
	"github.com/tin3ga/shortly/internal/database"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App, queries *database.Queries, ctx context.Context, rdb *redis.Client, ttl time.Duration, urlStr string, apiKey string) {
	app.Get("/", handler.Ping)

	// shortly
	app.Get("/api/v1/", func(c *fiber.Ctx) error {
		return handler.GetLinks(c, queries, ctx)
	})
	app.Get("/api/v1/:link", func(c *fiber.Ctx) error {
		return handler.GetLink(c, queries, ctx, rdb, ttl)
	})

	app.Post("/api/v1/shorten", func(c *fiber.Ctx) error {
		return handler.ShortenLink(c, queries, ctx, urlStr, apiKey)
	})
	app.Delete("/api/v1/shorten", func(c *fiber.Ctx) error {
		return handler.DeleteLink(c, queries, ctx)
	})
}
