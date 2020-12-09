package handler

import (
	"bytes"
	"fmt"
	"home/jonganebski/github/medium-rare/config"
	"image"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
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

	f, err := file.Open()
	imageSrc, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(500)
	}
	resizedImg := imaging.Resize(imageSrc, 800, 0, imaging.Lanczos)

	uuidWithHypen := uuid.New()
	uuid := strings.Replace(uuidWithHypen.String(), "-", "", -1)

	// ------
	// AWS S3
	// ------

	bucketName := config.Config("BUCKET_NAME")

	sess := connectAws()
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

	// ------
	// Local
	// ------

	// localURL := fmt.Sprintf("/image/%v", uuid+file.Filename)

	// if err = imaging.Save(resizedImg, "."+localURL); err != nil {
	// 	fmt.Println(err)
	// 	c.SendStatus(500)
	// }

	output := new(uploadPhotoByFileOutput)
	output.Success = 1
	// output.File.URL = "http://localhost:4000" + localURL
	output.File.URL = up.Location

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

func connectAws() *session.Session {
	AccessKeyID := config.Config("AWS_ACCESS_KEY_ID")
	SecretAccessKey := config.Config("AWS_SECRET_ACCESS_KEY")
	MyRegion := config.Config("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(MyRegion),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}

	return sess
}
