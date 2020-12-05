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
	app.Get("/read/:storyId", handler.ReadStory)

	app.Get("/blocks/:storyId", handler.ProvideStoryBlocks)

	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)

	app.Post("/upload/photo/byfile", middleware.Protected, handler.UploadPhotoByFilename)
	app.Post("/upload/story", middleware.Protected, handler.AddStory)
}
