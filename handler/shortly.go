package handler

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"

	"github.com/tin3ga/shortly/internal/database"
	"github.com/tin3ga/shortly/utils"
)

// Shorten Link model info
//
//	@Description	Shorten link Model
//	@Description	Url, Custom_alias
type ShortenLinkModel struct {
	Url          string `json:"url"`
	Custom_alias string `json:"custom_alias"`
	// expiration_date string `json:"expiration_date"`

}

// Delete Link model info
//
//	@Description	Delete Link Model
//	@Description	Url
type DeleteLinkModel struct {
	Url string `json:"url"`
}

// getLinks Fetch all links
//
//	@Summary		Fetch all links
//	@Description	Returns all links
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/links/all [get]
func GetLinks(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	data, err := queries.GetLinks(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cannot fetch links"})
	}
	log.Println("Fetching links")

	return c.JSON(data)

}

// getUserLinks Fetch all links associated to a user
//
//	@Summary		Fetch all user links
//	@Description	Returns all user links
//	@Tags			protected
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/links/userlinks [get]
func GetUserLinks(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {

	authHeader := c.Get("Authorization")

	userIDString, err := GetUserIDFromClaims(c, authHeader)
	if err != nil {
		log.Print(err)
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error parsing UserID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid UserID format"})
	}

	data, err := queries.GetUserLinks(ctx, userID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cannot fetch links"})
	}

	log.Println("Fetching User links")

	return c.Status(fiber.StatusOK).JSON(data)

}

// getLink Fetch a Original URL by Short URL
//
//	@Summary		Fetch a Original URL by Short URL
//	@Description	Redirects to the original URL
//	@Param			link	path	string	true	"Redirects to Original URL"
//	@Success		301
//	@Failure		404
//	@Router			/{link} [get]
func GetLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context, rdb *redis.Client, ttl time.Duration) error {
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

	// increment click count by one

	if err := queries.IncrementClickCount(ctx, link); err != nil {
		log.Print(err)
	}

	log.Println("Redirecting to: ", data.LongLink)
	return c.Redirect(data.LongLink, fiber.StatusMovedPermanently)

}

// shortenLink Insert an entry for a Short URL and Long URL
//
//	@Summary		Insert an entry for a Short URL and Long URL
//	@Description	Returns a Short URL
//	@Param			shorten_link	body	ShortenLinkModel	true	"Shorten a Link (custom alias is optional)"
//	@Tags			protected
//	@Security		BearerAuth
//	@Success		200
//	@Failure		400
//	@Failure		403
//	@Failure		500
//	@Router			/api/v1/links/shorten [post]
func ShortenLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context, apiKey string) error {
	url := new(ShortenLinkModel)

	if err := c.BodyParser(url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check if the URL starts with "https"
	if !strings.HasPrefix(url.Url, "https://") {
		log.Printf("Invalid URL scheme: %v", url.Url)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL must start with https://",
			"url":   url.Url,
		})

	}

	// Validate the URL using the external API
	result, err := utils.URLValidation(url.Url, apiKey)
	if err != nil {
		if err.Error() == "API error" {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Check that url is valid / try again later:)", "url": url.Url})
		}
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "External API error, try again later:)"})

	}

	// Handle malicious URL detection
	if result == "malicious" {
		log.Printf("Malicious url detected: %v", url.Url)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Url is malicious", "url": url.Url})

	} else {
		log.Printf("URL provided: %v is %v", url.Url, result)
	}

	var ShortLink string

	if url.Custom_alias == "" {
		ShortLink = randstr.Hex(8) // Generate a random 8 character string
	} else {
		ShortLink = url.Custom_alias
	}

	LongLink := url.Url
	uuidUser := uuid.New()

	authHeader := c.Get("Authorization")

	// type UserClaims struct {
	// 	Userid   string `json:"userid"`
	// 	Username string `json:"username"`
	// 	jwt.RegisteredClaims
	// }

	// cfg := config.InitializeConfig()
	// var jwtSecret = []byte(cfg.JWTSecret)

	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// // Parse and verify the token
	// token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	// Validate the algorithm
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return jwtSecret, nil
	// })

	// if err != nil || !token.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	// }

	// // Extract user claims
	// claims, ok := token.Claims.(*UserClaims)
	// if !ok {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token claims"})
	// }

	// // Store the user info in the context
	// c.Locals("userid", claims.Userid)
	// c.Locals("username", claims.Username)

	userIDString, err := GetUserIDFromClaims(c, authHeader)
	if err != nil {
		log.Print(err)
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error parsing UserID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid UserID format"})
	}
	log.Print(userID)

	params := database.CreateShortLinkParams{
		ID:        uuidUser,
		UserID:    userID,
		ShortLink: ShortLink,
		LongLink:  LongLink,
	}
	_, err = queries.CreateShortLink(ctx, params)
	if err != nil {
		log.Print(err)
		if err.Error() == "pq: duplicate key value violates unique constraint \"unique_short_link\"" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Duplicate short link, create a new alias"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create short link"})
	}
	log.Println("Created a shortened link: ", ShortLink)

	return c.JSON(fiber.Map{"Success": "Shortened link created", "url link": ShortLink})
}

// deleteLink Delete url data by short url
//
//	@Summary		Delete url data by short url
//	@Description	Returns a success message
//	@Param			url	body	DeleteLinkModel	true	"Delete a Link"
//	@Tags			protected
//	@Security		BearerAuth
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/api/v1/links/shorten [delete]
func DeleteLink(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	url := new(DeleteLinkModel)

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
