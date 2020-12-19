package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/photo"
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
func StoryRouter(api fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service, photoService photo.Service) {
	api.Get("/blocks/:storyId", provideStoryBlocks(storyService))
	api.Get("/recent-stories/:timestamp", provideRecentStories(userService, storyService))
	api.Post("/story", middleware.APIGuard, addStory(userService, storyService))
	api.Patch("/toggle-publish/:storyId/:toggle", middleware.APIGuard, togglePublish(userService, storyService))
	api.Patch("/toggle-like/:storyId", middleware.APIGuard, handleLikeCount(userService, storyService))
	api.Patch("/story/:storyId", middleware.APIGuard, editStory(storyService))
	api.Delete("/story/:storyId", middleware.APIGuard, removeStory(userService, storyService, commentService, photoService))
}

func provideRecentStories(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Params("timestamp")
		timestamp, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		stories, err := storyService.FindRecentStories(timestamp)
		if err != nil {
			return c.Status(404).SendString("Stories not found")
		}

		storyCards, err := composeStoryCardOutput(*stories, userService)
		if err != nil {
			return c.SendStatus(500)
		}
		fmt.Println(storyCards)

		return c.Status(200).JSON(storyCards)
	}
}

func togglePublish(userService user.Service, storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		toggleParam := c.Params("toggle")
		toggle, err := strconv.Atoi(toggleParam)
		if err != nil {
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
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}
		if story.CreatorID != userOID {
			return c.Status(403).SendString("You are not authorized.")
		}
		var makePublish bool = true
		if toggle < 0 {
			makePublish = false
		}
		err = storyService.PublishUnpublish(storyOID, makePublish)
		if err != nil {
			return c.Status(500).SendString("Failed to publish/unpublish")
		}
		return c.Status(200).JSON(fiber.Map{"isPublished": makePublish})
	}
}

func removeStory(userService user.Service, storyService story.Service, commentService comment.Service, photoService photo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
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

		// --- find story ---
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		// --- check current user is the author of the story ---
		if story.CreatorID != userOID {
			return c.Status(400).SendString("You are not authorized.")
		}

		// --- pull storyID from currnet user's storyIDs field ---
		err = userService.RemoveStoryID(userOID, storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- pull storyID from other users' likedStoryIDs field ---
		err = userService.RemoveManyLikedStoryIDs(storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- pull storyID from other users' savedStoryIDs field ---
		err = userService.RemoveManySavedStoryIDs(storyOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- pull every commentID of the story from other users' commentIDs field ---
		for _, commentID := range *story.CommentIDs {
			err = userService.RemoveManyUsersCommentID(commentID)
			if err != nil {
				return c.Status(500).SendString("Failed to delete")
			}
		}

		// --- remove all comment documents of the story ---
		err = commentService.RemoveComments(story.CommentIDs)
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
		if len(objects) != 0 {
			_, err = photoService.DeleteImagesOfS3(objects)
			if err != nil {
				fmt.Println(err)
				return c.SendStatus(500)
			}
		}

		// --- remove story document ---
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

		// --- find story ---
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		// --- check current user is the author of the story ---
		if userOID != story.CreatorID {
			return c.Status(400).SendString("You are not authorized")
		}

		// --- update story's blocks field ---
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

		// --- fill the fields need to be filled ---
		story.ID = ""
		story.CreatedAt = time.Now().Unix()
		story.UpdatedAt = time.Now().Unix()
		story.CreatorID = userOID
		story.LikedUserIDs = &[]primitive.ObjectID{}
		story.CommentIDs = &[]primitive.ObjectID{}
		story.ViewCount = 0
		story.EditorsPick = false
		// story.IsPublished = true
		story.IsPublished = false

		// --- insert story document ---
		storyOID, err := storyService.CreateStory(story)
		if err != nil {
			return c.Status(500).SendString("Failed to publish story")
		}

		// --- add the storyID into current user's storyIDs field ---
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
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			return c.SendStatus(500)
		}

		// --- find current user ---
		currentUser, err := userService.FindUserByID(userOID)
		if err != nil {
			return c.Status(404).SendString("User not found")
		}

		var isLikeExists bool = false
		for _, likedStoryID := range *currentUser.LikedStoryIDs {
			if likedStoryID == storyOID {
				isLikeExists = true
				break
			}
		}

		// --- determine update operator of mongodb ---
		var key string
		if !isLikeExists {
			key = "$push"
		}
		if isLikeExists {
			key = "$pull"
		}

		// --- manipulate story's LikedUserIDs field ---
		err = storyService.UpdateLikedUserIDs(storyOID, userOID, key)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- manipulate user's LikedStoryIDs field ---
		userService.UpdateLikedStoryIDs(userOID, storyOID, key)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		if !isLikeExists {
			return c.Status(200).SendString("1")
		}
		return c.Status(200).SendString("-1")
	}
}

func provideStoryBlocks(storyService story.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(500)
		}

		// --- find story ---
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		// return blocks of the story
		return c.JSON(story.Blocks)
	}
}
