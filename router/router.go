package router

import (
	"home/jonganebski/github/fibersteps-server/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.Home)

	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)
}
