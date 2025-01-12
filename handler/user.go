package handler

import (
	"context"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tin3ga/shortly/internal/database"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateUser(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	new_user := new(User)

	if err := c.BodyParser(new_user); err != nil {
		log.Print(err)
	}

	if new_user.Username == "" || new_user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and email required"})
	}
	// Email validation using regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(new_user.Email) {
		log.Print("Invalid email format:", new_user.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if new_user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password is required"})
	}
	// Password validation (simple example, could be extended for strength checks)
	if len(new_user.Password) < 8 {
		log.Print("Password is too short")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}

	hash, err := hashPassword(new_user.Password)
	if err != nil {

		log.Print(err)
	}
	uuid := uuid.New()

	createUserParams := database.CreateUserParams{
		ID:           uuid,
		Username:     new_user.Username,
		Email:        new_user.Email,
		PasswordHash: hash,
	}

	_, err = queries.CreateUser(ctx, createUserParams)
	if err != nil {
		// Handle duplicate username error
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A user with this username already exists"})
		}
		// Handle duplicate email error
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			log.Print("Duplicate email error:", err)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A user with this email already exists"})
		}
		// Log and return a generic error if it's not a duplicate key error
		log.Print("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})

	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": "User Created"})

}

func DeleteUser(c *fiber.Ctx, queries *database.Queries, ctx context.Context) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	input := new(PasswordInput)
	if err := c.BodyParser(input); err != nil {
		log.Print(err)
	}

	if input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password required"})
	}

	authHeader := c.Get("Authorization")
	userIDString, err := GetUserIDFromClaims(c, authHeader)
	if err != nil {
		log.Print(err)
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error parsing UserID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UserID format"})
	}

	userModel, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		log.Print(err)
	}

	if !CheckPasswordHash(userModel.PasswordHash, input.Password) {
		return c.JSON(fiber.Map{"error": "Unauthorized, invalid password"})
	}

	if err := queries.DeleteUser(ctx, userID); err != nil {
		log.Print(err)
		return c.JSON(fiber.Map{"error": "cannot delete user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": "user deleted"})

}
