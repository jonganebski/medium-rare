package middleware

import (
	"fmt"
	"home/jonganebski/github/medium-rare/config"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func getTokenFromCookie(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := c.Cookies("jwt")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(config.Config("JWT_SECRET")), nil
	})

	return token, err
}

// Protected only accepts logged-in users
func Protected(c *fiber.Ctx) error {
	userID := c.Locals("userId")
	username := c.Locals("username")
	if userID == nil || username == nil {
		return c.Redirect("/")
	}
	return c.Next()
}

// OnlyPublic only accepts loged-out users
// func OnlyPublic(c *fiber.Ctx) error {
// 	userID := c.Locals("userId")
// 	username := c.Locals("username")
// 	if userID != nil && username != nil {
// 		return c.Redirect(fmt.Sprintf("/%v", username))
// 	}
// 	return c.Next()
// }

// SpreadLocals gives local variables to context
func SpreadLocals(c *fiber.Ctx) error {
	token, err := getTokenFromCookie(c)

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["userId"]
			username := claims["username"]
			c.Locals("userId", userID)
			c.Locals("username", username)
		}
	}
	return c.Next()
}
