package handler

import "github.com/gofiber/fiber/v2"

// CreateUser creates a user
func CreateUser(c *fiber.Ctx) error {
	return c.SendStatus(201)
}

// Signin verify user password and gives jwt token
func Signin(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
