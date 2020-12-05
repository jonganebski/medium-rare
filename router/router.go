package router

import (
	"fmt"
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

	app.Post("/upload/photo/byfile", func(c *fiber.Ctx) error {

		type fileDetail struct {
			URL string `json:"url"`
		}

		type uploadPhotoByFileOutput struct {
			Success uint8      `json:"success"`
			File    fileDetail `json:"file"`
		}

		file, err := c.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}

		localURL := fmt.Sprintf("/image/%v", file.Filename)

		if err = c.SaveFile(file, "."+localURL); err != nil {
			fmt.Println(err)
		}

		output := new(uploadPhotoByFileOutput)
		output.Success = 1
		output.File.URL = "http://localhost:4000" + localURL

		return c.Status(200).JSON(output)
	})
}
