package router

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/tin3ga/shortly/handler"
	"github.com/tin3ga/shortly/internal/database"
	"github.com/tin3ga/shortly/middleware"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App, queries *database.Queries, ctx context.Context, rdb *redis.Client, ttl time.Duration, urlStr string, apiKey string, jwtsecret string) {
	app.Get("/", handler.Ping)
	app.Get("/:link", func(c *fiber.Ctx) error {
		return handler.GetLink(c, queries, ctx, rdb, ttl)
	})

	api := app.Group("api/v1")

	// user
	user := api.Group("/users")
	user.Post("/", func(c *fiber.Ctx) error {
		return handler.CreateUser(c, queries, ctx)

	})
	user.Delete("/delete", func(c *fiber.Ctx) error {
		return handler.DeleteUser(c, queries, ctx)
	})

	// auth
	auth := api.Group("/auth")
	auth.Post("/", func(c *fiber.Ctx) error {
		return handler.Login(c, queries, ctx, jwtsecret)
	})

	// shortly

	links := api.Group("links")
	links.Get("/all", func(c *fiber.Ctx) error {
		return handler.GetLinks(c, queries, ctx)
	})
	links.Get("/userlinks", func(c *fiber.Ctx) error {
		return handler.GetUserLinks(c, queries, ctx)
	})

	links.Post("/shorten", middleware.Protected(), func(c *fiber.Ctx) error {
		return handler.ShortenLink(c, queries, ctx, urlStr, apiKey)
	})
	links.Delete("/shorten", middleware.Protected(), func(c *fiber.Ctx) error {
		return handler.DeleteLink(c, queries, ctx)
	})
}
