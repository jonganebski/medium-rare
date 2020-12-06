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
	app.Get("/signout", middleware.Protected, handler.Signout)

	publicAPI := app.Group("/api")
	publicAPI.Get("/blocks/:storyId", handler.ProvideStoryBlocks)

	privateAPI := app.Group("/api", middleware.Protected)
	privateAPI.Post("/photo/byfile", handler.UploadPhotoByFilename)
	privateAPI.Delete("/photo", handler.DeletePhoto)
	privateAPI.Post("/story", handler.AddStory)
	privateAPI.Patch("/story/:storyId", handler.UpdateStory)
}
