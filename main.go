package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/redis/go-redis/v9"

	"github.com/tin3ga/shortly/cache"
	"github.com/tin3ga/shortly/config"
	"github.com/tin3ga/shortly/db"
	"github.com/tin3ga/shortly/router"

	"github.com/tin3ga/shortly/internal/database"

	_ "github.com/lib/pq"
	_ "github.com/tin3ga/shortly/docs"
)

const version = "0.4.0"

//	@title			Shortly API
//	@version		0.4.0
//	@description	This is a URL Shortener backend API built with Go.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.email	tinegagideon@outlook.com
//	@license.name	MIT License
//	@license.url	https://mit-license.org/
//	@host			shortly-5p7d.onrender.com
//	@BasePath		/

func main() {
	cfg := config.InitializeConfig()

	log.Printf("Starting server on port %v", cfg.Port)
	log.Printf("Serving version %v", version)

	db, err := db.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	queries := database.New(db)

	// Caching - Redis setup

	var rdb *redis.Client

	if cfg.EnableCaching {
		rdb, err := cache.InitializeRedis(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

		if err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}
		defer rdb.Close() // Close Redis client when no longer needed
		log.Printf("--Cache TTL: %v", cfg.CacheTTL)

	}

	// end Redis setup

	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Rate limiter

	// Set up in-memory store for the rate limiter

	if cfg.EnableRateLimiting {
		limiterCfg := limiter.Config{
			Max:        cfg.MaxConnectionsLimit,
			Expiration: cfg.RateLimitExpiration,
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Request limit reached! Try again later:)"})
			},
			SkipFailedRequests:     cfg.SkipFailedRequests,
			SkipSuccessfulRequests: cfg.SkipSuccessfulRequests,
		}
		app.Use(limiter.New(limiterCfg))

		log.Printf("Rate Limiting Enabled: %v", cfg.EnableRateLimiting)
		log.Printf("--Max Connections Limit: %v", cfg.MaxConnectionsLimit)
		log.Printf("--Rate Limit Expiration: %v", cfg.RateLimitExpiration)
		log.Printf("--Skip Failed Requests: %v", cfg.SkipFailedRequests)
		log.Printf("--Skip Successful Requests: %v", cfg.SkipSuccessfulRequests)

	}

	// health check
	// Provide a minimal config
	app.Use(healthcheck.New())

	// Initialize default config (Assign the middleware to /metrics)
	// app.Get("/metrics", monitor.New())
	metricsCfg := monitor.Config{
		Title:      cfg.Title,
		FontURL:    cfg.FontURL,
		ChartJsURL: "https://cdn.jsdelivr.net/npm/chart.js@2.9/dist/Chart.bundle.min.js",
		APIOnly:    false,
		Next:       nil,
	}

	app.Get("/metrics", monitor.New(metricsCfg))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	router.SetupRoutes(app, queries, ctx, rdb, cfg.CacheTTL, cfg.APIKey, cfg.JWTSecret)

	app.Listen(":" + cfg.Port)
}
