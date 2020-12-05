package router

import (
	"home/jonganebski/github/fibersteps-server/handler"
	"home/jonganebski/github/fibersteps-server/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.Home)
	app.Get("/new-story", middleware.Protected, handler.NewStory)

	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)

	app.Post("/upload/photo/byfile", handler.UploadPhotoByFilename)
	app.Post("/upload/story", handler.AddStory)
}
