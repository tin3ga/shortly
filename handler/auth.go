package handler

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tin3ga/shortly/config"
	"github.com/tin3ga/shortly/internal/database"
)

func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetUserIDFromClaims(c *fiber.Ctx, authHeader string) (string, error) {
	type UserClaims struct {
		Userid   string `json:"userid"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
	cfg := config.InitializeConfig()
	var jwtSecret = []byte(cfg.JWTSecret)

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("missing or invalid token")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
		return "", err
	}
	// Extract user claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", err
	}

	// Store the user info in the context
	c.Locals("userid", claims.Userid)
	c.Locals("username", claims.Username)

	return claims.Userid, nil

}

func GenerateToken(userModel database.User, jwtsecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = userModel.ID
	claims["username"] = userModel.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtsecret))
	if err != nil {
		return "", err
	}
	return t, nil

}

func ValidateToken(token *jwt.Token, username string) bool {
	claims := token.Claims.(jwt.MapClaims)
	uname := claims["username"]

	return username == uname
}

func Login(c *fiber.Ctx, queries *database.Queries, ctx context.Context, jwtsecret string) error {
	type UserInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	login_identity := new(UserInput)

	if err := c.BodyParser(login_identity); err != nil {
		log.Print(err)
	}

	var userModel database.User
	var err error
	if isEmail(login_identity.Identity) {
		userModel, err = queries.GetUserByEmail(ctx, login_identity.Identity)
	} else {
		userModel, err = queries.GetUserByUsername(ctx, login_identity.Identity)

	}
	log.Print(userModel)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(fiber.Map{"error": "Unauthorized, user does not exist"})

		}

		log.Print(err)
	}

	// Invalidate user with  wrong password
	if !CheckPasswordHash(userModel.PasswordHash, login_identity.Password) {
		return c.JSON(fiber.Map{"error": "Unauthorized, invalid password"})

	}

	token, err := GenerateToken(userModel, jwtsecret)
	if err != nil {
		log.Printf("failed to generate token %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate token"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": "Authorized", "Data": token})

}
