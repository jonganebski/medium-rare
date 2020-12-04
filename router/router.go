package router

import (
	"home/jonganebski/github/fibersteps-server/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("editor", fiber.Map{}, "layout/main")
	})

	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)
}
