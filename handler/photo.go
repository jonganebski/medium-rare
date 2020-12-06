package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

	uuidWithHypen := uuid.New()
	uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)

	localURL := fmt.Sprintf("/image/%v", uuid+file.Filename)

	if err = c.SaveFile(file, "."+localURL); err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}

	output := new(uploadPhotoByFileOutput)
	output.Success = 1
	output.File.URL = "http://localhost:4000" + localURL

	return c.Status(200).JSON(output)
}

// DeletePhoto removes photo
func DeletePhoto(c *fiber.Ctx) error {
	// For now editorjs does not provide eventlistener for deleting a photo while editing.
	// Maybe I can customize it but it will take some time and it's not that worth it.
	// It's also on github issue https://github.com/editor-js/image/issues/54
	// Wait for editorjs update.
	return c.SendStatus(200)
}
