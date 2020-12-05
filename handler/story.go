package handler

import (
	"fmt"
	"home/jonganebski/github/fibersteps-server/config"
	"home/jonganebski/github/fibersteps-server/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var storyCollectionName = config.Config("COLLECTION_STORY")

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

	storyCollection := mg.Db.Collection(storyCollectionName)
	story := new(model.Story)

	if err := c.BodyParser(story); err != nil {
		fmt.Println("error at body parser")
		return c.SendStatus(400)
	}

	userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
	if err != nil {
		fmt.Println("error at conversion")
		return c.SendStatus(500)
	}

	story.ID = ""
	story.CreatedAt = time.Now().Unix()
	story.UpdatedAt = time.Now().Unix()
	story.CreatorID = userOID
	story.LikedUserIDs = &[]primitive.ObjectID{}
	story.CommentIDs = &[]primitive.ObjectID{}
	story.ViewCount = 0

	fmt.Println(story)

	_, err = storyCollection.InsertOne(c.Context(), story)
	if err != nil {
		fmt.Println("error at insertion")
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(story)
}
