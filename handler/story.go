package handler

import "github.com/gofiber/fiber/v2"

// Home renders homepage
func Home(c *fiber.Ctx) error {

	return c.Render("editor", fiber.Map{"userId": c.Locals("userId")}, "layout/main")
}
