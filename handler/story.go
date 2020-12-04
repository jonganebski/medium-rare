package handler

import "github.com/gofiber/fiber/v2"

// Home renders homepage
func Home(c *fiber.Ctx) error {

	return c.Render("home", fiber.Map{"userId": c.Locals("userId")}, "layout/main")
}

// NewStory renders a pagr where a user writes a new story
func NewStory(c *fiber.Ctx) error {
	return c.Render("newStory", fiber.Map{"userId": c.Locals("userId")}, "layout/main")
}
