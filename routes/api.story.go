package routes

import (
	"fmt"
	myaws "home/jonganebski/github/medium-rare/aws"
	"home/jonganebski/github/medium-rare/config"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StoryRouter has api routes for stories
func StoryRouter(app fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service) {
	app.Get("/blocks/:storyId", provideStoryBlocks(storyService))
	app.Post("/story", middleware.APIGuard, addStory(userService, storyService))
	app.Post("/like/:storyId/:plusMinus", middleware.APIGuard, handleLikeCount(userService, storyService))
	app.Patch("/story/:storyId", middleware.APIGuard, editStory(storyService))
	app.Delete("/story/:storyId", middleware.APIGuard, removeStory(userService, storyService, commentService))
}

func removeStory(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyID")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}
		if story.CreatorID != userOID {
			return c.Status(400).SendString("You are not authorized.")
		}
		err = userService.RemoveStoryID(userOID, storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		err = userService.RemoveManyLikedStoryIDs(storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		err = userService.RemoveManySavedStoryIDs(storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		err = commentService.RemoveComments(story.CommentIDs)
		if err != nil {
			return c.Status(500).SendString("Failed to delete")
		}
		err = userService.RemoveManyCommentIDs(story.CommentIDs)
		if err != nil {
			return c.Status(500).SendString("Failed to delete")
		}

		// --- remove related images in AWS S3 ---

		objects := make([]*s3.ObjectIdentifier, 0)

		for _, block := range story.Blocks {
			if block.Type == "image" {
				fileName := strings.Split(block.Data.File.URL, "amazonaws.com/")[1]
				objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(fileName)})
			}
		}

		sess := myaws.ConnectAws()
		svc := s3.New(sess)
		bucketName := config.Config("BUCKET_NAME")
		_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{Bucket: aws.String(bucketName), Delete: &s3.Delete{Objects: objects, Quiet: aws.Bool(true)}})
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}

		err = storyService.RemoveStory(storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to delete")
		}
		return c.SendStatus(204)
	}
}

func editStory(storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		editedStory := new(model.Story)

		if err := c.BodyParser(editedStory); err != nil {
			fmt.Println("error at body parser")
			return c.SendStatus(400)
		}

		storyID := c.Params("storyId")

		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			fmt.Println("error at conversion")
			return c.SendStatus(500)
		}
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			fmt.Println("error at conversion")
			return c.SendStatus(500)
		}

		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		if userOID != story.CreatorID {
			return c.Status(400).SendString("You are not authorized")
		}

		err = storyService.UpdateStoryBlock(storyOID, &editedStory.Blocks)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		return c.SendStatus(200)
	}
}

func addStory(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		story := new(model.Story)

		if err := c.BodyParser(story); err != nil {
			return c.SendStatus(400)
		}

		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		story.ID = ""
		story.CreatedAt = time.Now().Unix()
		story.UpdatedAt = time.Now().Unix()
		story.CreatorID = userOID
		story.LikedUserIDs = &[]primitive.ObjectID{}
		story.CommentIDs = &[]primitive.ObjectID{}
		story.ViewCount = 0
		story.EditorsPick = false
		story.IsPublished = true

		storyOID, err := storyService.CreateStory(story)
		if err != nil {
			return c.Status(500).SendString("Failed to publish story")
		}

		err = userService.AddStoryID(userOID, *storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		return c.Status(201).SendString(storyOID.Hex())
	}
}

func handleLikeCount(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		p := c.Params("plusMinus") // bigger than zero -> increase like / smaller than zero -> decrease like
		plusMinus, err := strconv.Atoi(p)
		if err != nil || plusMinus == 0 {
			return c.SendStatus(400)
		}
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}
		// --- check the user is allowed to like or unlike the story ---
		var isAllowed bool
		if plusMinus > 0 {
			isAllowed = true
			for _, likedStoryID := range *currentUser.LikedStoryIDs {
				if likedStoryID == storyOID {
					isAllowed = false
				}
			}
		}
		if plusMinus < 0 {
			isAllowed = false
			for _, likedStoryID := range *currentUser.LikedStoryIDs {
				if likedStoryID == storyOID {
					isAllowed = true
				}
			}
		}
		if !isAllowed {
			return c.Status(400).SendString("You can't like or unlike twice.")
		}
		var key string
		if plusMinus > 0 {
			key = "$push"
		}
		if plusMinus < 0 {
			key = "$pull"
		}
		err = storyService.UpdateLikedUserIDs(storyOID, userOID, key)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		userService.UpdateLikedStoryIDs(userOID, storyOID, key)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}
		return c.SendStatus(200)
	}
}

func provideStoryBlocks(storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		return c.JSON(story.Blocks)
	}
}
