package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Verify Connection
//
//	@Summary		Checks connectivity
//	@Description	Returns pong
//	@Success		200
//	@Router			/ [get]
func Ping(c *fiber.Ctx) error {
	log.Println("Ping.....")
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
