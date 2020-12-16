package routes

import (
	"fmt"
	"home/jonganebski/github/medium-rare/middleware"
	"home/jonganebski/github/medium-rare/model"
	"home/jonganebski/github/medium-rare/package/comment"
	"home/jonganebski/github/medium-rare/package/story"
	"home/jonganebski/github/medium-rare/package/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentOutput struct {
	CommentID    string `json:"commentId"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatarUrl"`
	CreatedAt    int64  `json:"createdAt"`
	Text         string `json:"text"`
	IsAuthorized bool   `json:"isAuthorized"`
}

// CommentRouter has api routes for comment
func CommentRouter(api fiber.Router, userService user.Service, storyService story.Service, commentService comment.Service) {
	api.Get("/comment/:storyId", provideComments(userService, storyService, commentService))
	api.Post("/comment/:storyId", middleware.APIGuard, addComment(userService, storyService, commentService))
	api.Delete("/comment/:commentId", middleware.APIGuard, removeComment(userService, storyService, commentService))
}

func removeComment(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		commentID := c.Params("commentId")
		commentOID, err := primitive.ObjectIDFromHex(commentID)
		if err != nil {
			c.SendStatus(500)
		}
		userOID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", c.Locals("userId")))
		if err != nil {
			c.SendStatus(500)
		}

		// --- find comment ---
		comment, err := commentService.FindComment(commentOID)
		if err != nil {
			return c.Status(404).SendString("Comment not found")
		}

		// --- check current user is who wrote the comment ---
		if comment.CreatorID != userOID {
			return c.Status(403).SendString("You are not authorized.")
		}

		// --- remove commentID from current user's commentIDs field ---
		err = userService.RemoveCommentID(userOID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- remove commentID from story's commentIDs field ---
		err = storyService.RemoveCommentID(comment.StoryID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Failed to update")
		}

		// --- remove comment document itself ---
		err = commentService.RemoveComment(commentOID)
		if err != nil {
			return c.Status(500).SendString("Failed to delete")
		}

		return c.SendStatus(200)
	}
}

func addComment(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		comment := new(model.Comment)
		if err := c.BodyParser(comment); err != nil {
			return c.SendStatus(400)
		}
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

		// --- fill the fields ---
		comment.ID = ""
		comment.CreatorID = userOID
		comment.CreatedAt = time.Now().Unix()
		comment.StoryID = storyOID

		// --- insert new comment document ---
		comment, err = commentService.CreateComment(comment)
		if err != nil {
			return c.Status(500).SendString("Create comment failed")
		}
		commentOID, err := primitive.ObjectIDFromHex(comment.ID)
		if err != nil {
			return c.SendStatus(500)
		}

		// --- push commentID into current user's commentIDs field ---
		err = userService.AddCommentID(userOID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Update failed")
		}

		// --- push commentID into story's commentIDs field ---
		err = storyService.AddCommentID(storyOID, commentOID)
		if err != nil {
			return c.Status(500).SendString("Update failed")
		}

		// --- make a comment output ---
		output := new(commentOutput)
		output.CommentID = comment.ID
		output.AvatarURL = currentUser.AvatarURL
		output.CreatedAt = comment.CreatedAt
		output.Text = comment.Text
		output.Username = currentUser.Username
		output.IsAuthorized = true

		return c.Status(201).JSON(output)
	}
}

func provideComments(userService user.Service, storyService story.Service, commentService comment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		storyID := c.Params("storyId")
		storyOID, err := primitive.ObjectIDFromHex(storyID)
		if err != nil {
			return c.SendStatus(400)
		}

		// --- find story ---
		story, err := storyService.FindStoryByID(storyOID)
		if err != nil {
			return c.Status(404).SendString("Story not found")
		}

		// --- find comments of the story ---
		comments, err := commentService.FindComments(story.CommentIDs)
		if err != nil {
			return c.Status(404).SendString("Comments not found")
		}

		// --- make slice of comment outputs ---
		outputItem := new(commentOutput)
		output := make([]commentOutput, 0)
		for _, comment := range *comments {
			// find creator if a comment
			creator, err := userService.FindUserByID(comment.CreatorID)
			if err != nil {
				return c.Status(404).SendString("Commenter not found")
			}
			// make a comment output and append to the slice
			outputItem.CommentID = comment.ID
			outputItem.Username = creator.Username
			outputItem.AvatarURL = creator.AvatarURL
			outputItem.CreatedAt = comment.CreatedAt
			outputItem.Text = comment.Text
			outputItem.IsAuthorized = (creator.ID == c.Locals("userId"))
			output = append(output, *outputItem)
		}

		return c.Status(200).JSON(output)
	}
}
