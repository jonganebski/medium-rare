package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Home renders homepage
func Home(c *fiber.Ctx) error {

	return c.Render("home", fiber.Map{"path": c.Path(), "userId": c.Locals("userId")}, "layout/main")
}

// NewStory renders a pagr where a user writes a new story
func NewStory(c *fiber.Ctx) error {
	return c.Render("newStory", fiber.Map{"path": c.Path(), "userId": c.Locals("userId")}, "layout/main")
}

// UploadPhotoByFilename saves photo that user attached on the story 'when attatchment occurs'
func UploadPhotoByFilename(c *fiber.Ctx) error {
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
		return c.SendStatus(400)
	}

	localURL := fmt.Sprintf("/image/%v", file.Filename)

	if err = c.SaveFile(file, "."+localURL); err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}

	output := new(uploadPhotoByFileOutput)
	output.Success = 1
	output.File.URL = "http://localhost:4000" + localURL

	return c.Status(200).JSON(output)
}

// AddStory creates a new story
func AddStory(c *fiber.Ctx) error {
	return c.SendStatus(201)
}
