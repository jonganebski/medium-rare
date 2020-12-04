package router

import (
	"fmt"
	"home/jonganebski/github/fibersteps-server/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("editor", fiber.Map{}, "layout/main")
	})
	app.Post("/", func(c *fiber.Ctx) error {
		type input struct {
			Data string `json:"data"`
		}

		data := new(input)

		c.BodyParser(data)

		fmt.Println(data)

		return c.SendStatus(200)
	})
	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)
}
