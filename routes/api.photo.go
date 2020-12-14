package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/package/photo"
	"image"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// --- structure of the output editorjs wants ---
type fileDetail struct {
	URL string `json:"url"`
}
type uploadPhotoByFileOutput struct {
	Success uint8      `json:"success"`
	File    fileDetail `json:"file"`
}

// ImageRouter has routes related uploading images
func ImageRouter(api fiber.Router, photoService photo.Service) {
	api.Post("/photo/byfile", middleware.APIGuard, uploadPhotoByFilename(photoService))
	api.Delete("/photos", middleware.APIGuard, deletePhotos(photoService))
}

func deletePhotos(photoService photo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type deletePhotosInput struct {
			Images []string `json:"images"`
		}
		images := new(deletePhotosInput)
		if err := c.BodyParser(images); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}
		objects := make([]*s3.ObjectIdentifier, 0)
		for _, url := range images.Images {
			fileName := strings.Split(url, "amazonaws.com/")[1]
			objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(fileName)})

		}
		if len(objects) != 0 {
			_, err := photoService.DeleteImagesOfS3(objects)
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
		}

		return c.SendStatus(204)
	}
}

func uploadPhotoByFilename(photoService photo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		output := new(uploadPhotoByFileOutput)
		file, err := c.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}

		// --- open & decode image file ---
		f, err := file.Open()
		imageSrc, _, err := image.Decode(f)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		// --- resize the image ---
		resizedImg := imaging.Resize(imageSrc, 1000, 0, imaging.Lanczos)

		// --- make filename with uuid ---
		uuidWithHypen := uuid.New()
		uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)
		filename := uuid + file.Filename

		up, err := photoService.UploadImageToS3(resizedImg, filename)
		if err != nil {
			fmt.Println("Failed to upload file")
			return c.SendStatus(500)
		}
		output.File.URL = up.Location
		output.Success = 1

		// ------
		// Local
		// ------

		// localURL := fmt.Sprintf("/image/%v", uuid+file.Filename)

		// if err = imaging.Save(resizedImg, "."+localURL); err != nil {
		// 	fmt.Println(err)
		// 	c.SendStatus(500)
		// }
		// output.File.URL = "http://localhost:4000" + localURL

		// output.Success = 1

		return c.Status(200).JSON(output)
	}
}
