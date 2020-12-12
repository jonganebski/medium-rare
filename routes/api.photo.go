package routes

import (
	"bytes"
	"fmt"
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/middleware"
	"image"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ImageRouter has routes related uploading images
func ImageRouter(app fiber.Router) {
	api := app.Group("/api")
	api.Post("/photo/byfile", middleware.APIGuard, uploadPhotoByFilename)
}

func uploadPhotoByFilename(c *fiber.Ctx) error {

	type fileDetail struct {
		URL string `json:"url"`
	}

	type uploadPhotoByFileOutput struct {
		Success uint8      `json:"success"`
		File    fileDetail `json:"file"`
	}

	output := new(uploadPhotoByFileOutput)

	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(400)
	}

	f, err := file.Open()
	imageSrc, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}
	resizedImg := imaging.Resize(imageSrc, 1000, 0, imaging.Lanczos)

	uuidWithHypen := uuid.New()
	uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)

	// ------
	// AWS S3
	// ------

	bucketName := config.Config("BUCKET_NAME")

	sess := myaws.ConnectAws()
	uploader := s3manager.NewUploader(sess)

	filename := uuid + file.Filename

	buf := new(bytes.Buffer)
	imaging.Encode(buf, resizedImg, imaging.JPEG)
	reader := bytes.NewReader(buf.Bytes())

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   reader,
	})

	if err != nil {
		fmt.Println("Failed to upload file")
		return c.SendStatus(500)
	}
	output.File.URL = up.Location

	// ------
	// Local
	// ------

	// localURL := fmt.Sprintf("/image/%v", uuid+file.Filename)

	// if err = imaging.Save(resizedImg, "."+localURL); err != nil {
	// 	fmt.Println(err)
	// 	c.SendStatus(500)
	// }
	// output.File.URL = "http://localhost:4000" + localURL

	output.Success = 1

	return c.Status(200).JSON(output)
}
