package router

import (
	"home/jonganebski/github/medium-rare/handler"
	"home/jonganebski/github/medium-rare/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", handler.Home)
	app.Get("/new-story", middleware.Protected, handler.NewStory)
	app.Get("/read/:storyId", handler.ReadStory)
	app.Get("/edit-story/:storyId", handler.EditStory)

	app.Post("/signup", handler.CreateUser)
	app.Post("/signin", handler.Signin)

	api := app.Group("/api", middleware.Protected)
	api.Get("/blocks/:storyId", handler.ProvideStoryBlocks)
	api.Post("/photo/byfile", handler.UploadPhotoByFilename)
	api.Post("/story", handler.AddStory)
	api.Patch("/story/:storyId", handler.UpdateStory)
}
